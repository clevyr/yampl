package processor

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gabe565.com/utils/coloryaml"
	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/util"
	"github.com/clevyr/yampl/internal/visitor"
	"gopkg.in/yaml.v3"
)

type FileTask struct {
	Path    string
	Mode    os.FileMode
	Content string
}

var ErrStdinInplace = errors.New("-i or --inplace may not be used with stdin")

func Process(conf *config.Config, args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return processStdin(conf, stdin, stdout)
	}
	return walkPaths(conf, args, stdout)
}

func processStdin(conf *config.Config, stdin io.Reader, stdout io.Writer) error {
	if conf.Inplace {
		return ErrStdinInplace
	}

	s, err := templateReader(conf, "stdin", stdin)
	if err != nil {
		return err
	}

	_, err = coloryaml.WriteString(stdout, s)
	return err
}

func walkPaths(conf *config.Config, args []string, stdout io.Writer) error {
	var hasDir bool
	for _, arg := range args {
		if stat, err := os.Lstat(arg); err == nil {
			if stat.IsDir() {
				hasDir = true
				break
			}
		}
	}

	multiFile := len(args) > 1 || hasDir
	if !conf.NoSourceComment {
		conf.NoSourceComment = !multiFile
	}

	tasks, err := collectTasks(conf, args, multiFile)
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if i != 0 && !conf.Inplace {
			if _, err := io.WriteString(stdout, "---\n"); err != nil {
				return err
			}
		}

		if err := FlushFile(conf, stdout, task); err != nil {
			return err
		}
	}

	return nil
}

func collectTasks(conf *config.Config, args []string, logErrors bool) ([]FileTask, error) {
	var errs []error
	var tasks []FileTask

	for _, arg := range args {
		if err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if logErrors {
					slog.Error("Failed to template file", "error", err)
				}
				errs = append(errs, err)
				return nil
			}

			if d.IsDir() || path != arg && !util.IsYaml(path) {
				return nil
			}

			task, err := TemplateFile(conf, path)
			if err != nil {
				if logErrors {
					slog.Error("Failed to template file", "path", path, "error", err)
				}
				errs = append(errs, err)
				return nil
			}
			tasks = append(tasks, *task)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}
	return tasks, nil
}

func TemplateFile(conf *config.Config, path string) (*FileTask, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	s, err := templateReader(conf, path, f)
	if err != nil {
		return nil, err
	}

	return &FileTask{
		Path:    path,
		Mode:    stat.Mode(),
		Content: s,
	}, nil
}

func FlushFile(conf *config.Config, w io.Writer, task FileTask) error {
	s := task.Content
	if !conf.Inplace {
		if !conf.NoSourceComment {
			source := "# Source: " + task.Path + "\n"
			if !strings.HasPrefix(s, "---") {
				s = source + s
			}
			if strings.Contains(s, "---") {
				s = strings.ReplaceAll(s, "---\n", "---\n"+source)
			}
		}

		_, err := coloryaml.WriteString(w, s)
		return err
	}

	temp, err := os.CreateTemp("", "yampl_*_"+filepath.Base(task.Path))
	if err != nil {
		return err
	}
	defer func() {
		_ = temp.Close()
		_ = os.Remove(temp.Name())
	}()

	if _, err := temp.WriteString(s); err != nil {
		return err
	}

	if err := temp.Chmod(task.Mode); err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	if err := os.Rename(temp.Name(), task.Path); err != nil {
		slog.Debug("Failed to rename file. Attempting to copy contents.",
			"from", temp.Name(),
			"to", task.Path,
			"error", err,
		)

		out, err := os.OpenFile(task.Path, os.O_WRONLY|os.O_TRUNC, task.Mode)
		if err != nil {
			return err
		}

		if _, err := out.WriteString(s); err != nil {
			return err
		}

		if err := out.Close(); err != nil {
			return err
		}
	}

	return nil
}

const indicator = "#_yampl_newline\n"

var searchRe = regexp.MustCompile(`\n\n(\s+)?`)

func templateReader(conf *config.Config, path string, r io.Reader) (string, error) {
	v := visitor.NewTemplateComments(conf, path)

	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	b = searchRe.ReplaceAll(b, []byte("\n$1"+indicator+"$1"))

	decoder := yaml.NewDecoder(bytes.NewReader(b))
	var buf strings.Builder
	buf.Grow(len(b))

	for {
		var n yaml.Node

		if err := decoder.Decode(&n); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}

		if buf.Len() > 0 {
			buf.WriteString("---\n")
		}

		if err := v.Run(&n); err != nil {
			return "", err
		}

		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(conf.Indent)
		if err := encoder.Encode(&n); err != nil {
			_ = encoder.Close()
			return "", err
		}

		if err := encoder.Close(); err != nil {
			return "", err
		}
	}

	var result strings.Builder
	result.Grow(buf.Len())
	for line := range strings.Lines(buf.String()) {
		if strings.HasSuffix(line, indicator) {
			result.WriteByte('\n')
		} else {
			result.WriteString(line)
		}
	}
	return result.String(), nil
}
