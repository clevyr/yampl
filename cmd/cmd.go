package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/coloryaml"
	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/util"
	"github.com/clevyr/yampl/internal/visitor"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

//nolint:gochecknoglobals
var description = `Yampl (yaml + tmpl) templates YAML values based on line-comments.
YAML data can be piped to stdin or files/dirs can be passed as arguments.

Full reference at ` + termenv.Hyperlink("https://github.com/clevyr/yampl#readme", "github.com/clevyr/yampl")

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "yampl [files | dirs] [-v key=value...]",
		Short:             "Inline YAML templating via line-comments",
		Long:              description,
		DisableAutoGenTag: true,
		ValidArgsFunction: validArgs,
		RunE:              run,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	conf.RegisterCompletions(cmd)
	visitor.RegisterCompletion(cmd)
	cmd.SetContext(config.WithContext(context.Background(), conf))

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
}

var ErrStdinInplace = errors.New("-i or --inplace may not be used with stdin")

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	if len(args) == 0 {
		if f, ok := cmd.InOrStdin().(*os.File); ok {
			if isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd()) {
				return cmd.Help()
			}
		}

		if conf.Inplace {
			return ErrStdinInplace
		}

		s, err := templateReader(conf, "stdin", cmd.InOrStdin())
		if err != nil {
			return err
		}

		if _, err := coloryaml.WriteString(cmd.OutOrStdout(), s); err != nil {
			return err
		}

		return nil
	}

	return walkPaths(cmd, conf, args)
}

func walkPaths(cmd *cobra.Command, conf *config.Config, args []string) error {
	var hasDir bool
	for _, arg := range args {
		if stat, err := os.Lstat(arg); err == nil {
			if stat.IsDir() {
				hasDir = true
				break
			}
		}
	}

	logErrors := len(args) > 1 || hasDir
	if !conf.NoSourceComment {
		conf.NoSourceComment = len(args) <= 1 && !hasDir
	}

	var printSeparator bool
	var errs []error
	for _, arg := range args {
		if err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if logErrors {
					slog.Error("Failed to template file", "error", err)
					printSeparator = true
				}
				errs = append(errs, err)
				return nil
			}

			if d.IsDir() || path != arg && !util.IsYaml(path) {
				return nil
			}

			if printSeparator && !conf.Inplace {
				printSeparator = false
				if _, err := io.WriteString(cmd.OutOrStdout(), "---\n"); err != nil {
					return err
				}
			}

			if err := openAndTemplateFile(conf, cmd.OutOrStdout(), path); err != nil {
				if logErrors {
					slog.Error("Failed to template file", "error", err)
				}
				errs = append(errs, err)
			}
			printSeparator = true
			return nil
		}); err != nil {
			return err
		}
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errors.Join(errs...)
	}
}

func openAndTemplateFile(conf *config.Config, w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	s, err := templateReader(conf, path, f)
	if err != nil {
		return err
	}

	_ = f.Close()

	if !conf.Inplace {
		if !conf.NoSourceComment {
			source := "# Source: " + path + "\n"
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

	temp, err := os.CreateTemp("", "yampl_*_"+filepath.Base(path))
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

	if err := temp.Chmod(stat.Mode()); err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	if err := os.Rename(temp.Name(), path); err != nil {
		slog.Debug("Failed to rename file. Attempting to copy contents.",
			"from", temp.Name(),
			"to", path,
			"error", err,
		)

		out, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, stat.Mode())
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
