package cmd

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/clevyr/yampl/internal/colorize"
	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/util"
	"github.com/clevyr/yampl/internal/visitor"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var version = "beta"

//nolint:gochecknoglobals
var description = `Yampl (yaml + tmpl) templates YAML values based on line-comments.
YAML data can be piped to stdin or files/dirs can be passed as arguments.

Full reference at ` + termenv.Hyperlink("https://github.com/clevyr/yampl#readme", "github.com/clevyr/yampl")

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "yampl [files | dirs] [-v key=value...]",
		Short:             "Inline YAML templating via line-comments",
		Long:              description,
		DisableAutoGenTag: true,
		ValidArgsFunction: validArgs,
		Version:           buildVersion(version),
		RunE:              run,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	cmd.InitDefaultVersionFlag()
	conf.RegisterCompletions(cmd)
	visitor.RegisterCompletion(cmd, conf)
	cmd.SetContext(config.WithContext(context.Background(), conf))
	return cmd
}

func validArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
}

var ErrStdinInplace = errors.New("-i or --inplace may not be used with stdin")

func run(cmd *cobra.Command, args []string) error {
	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("config missing from command context")
	}

	if err := conf.Load(cmd); err != nil {
		return err
	}

	if conf.Completion != "" {
		return completion(cmd, conf.Completion)
	}

	cmd.SilenceUsage = true

	if len(args) == 0 {
		if isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			return cmd.Help()
		}

		if conf.Inplace {
			return ErrStdinInplace
		}

		var size int64
		if stat, err := os.Stdin.Stat(); err == nil {
			size = stat.Size()
		}

		s, err := templateReader(conf, "stdin", os.Stdin, size)
		if err != nil {
			return err
		}

		if err := colorize.WriteString(cmd.OutOrStdout(), s); err != nil {
			return err
		}
	}

	withSrcComment := len(args) > 1
	if !withSrcComment {
		for _, arg := range args {
			stat, err := os.Stat(arg)
			if err != nil {
				return err
			}
			if withSrcComment = stat.IsDir(); withSrcComment {
				break
			}
		}
	}

	var i int
	for _, arg := range args {
		if err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !util.IsYaml(path) {
				return err
			}

			if !conf.Inplace && i != 0 {
				if _, err := io.WriteString(cmd.OutOrStdout(), "---\n"); err != nil {
					return err
				}
			}
			i++

			return openAndTemplateFile(conf, cmd.OutOrStdout(), arg, path, withSrcComment)
		}); err != nil {
			return err
		}
	}

	return nil
}

func openAndTemplateFile(conf *config.Config, w io.Writer, dir, path string, withSrcComment bool) error {
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

	s, err := templateReader(conf, path, f, stat.Size())
	if err != nil {
		return err
	}

	_ = f.Close()

	if conf.Inplace {
		temp, err := os.CreateTemp("", "yampl_*_"+filepath.Base(path))
		if err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(temp.Name())
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
			log.Trace().Msg("failed to rename file, attempting to copy contents")

			in, err := os.Open(temp.Name())
			if err != nil {
				return err
			}
			defer func() {
				_ = in.Close()
			}()

			out, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, stat.Mode())
			if err != nil {
				return err
			}

			if _, err := io.Copy(out, in); err != nil {
				return err
			}

			if err := out.Close(); err != nil {
				return err
			}
		}
	} else {
		if !conf.NoSourceComment {
			rel := path
			if !withSrcComment {
				if rel, err = filepath.Rel(dir, path); err == nil && rel != "." {
					withSrcComment = true
				}
			}
			if withSrcComment {
				source := "# Source: " + rel + "\n"
				if !strings.HasPrefix(s, "---") {
					s = source + s
				}
				if strings.Contains(s, "---") {
					s = strings.ReplaceAll(s, "---\n", "---\n"+source)
				}
			}
		}

		if err := colorize.WriteString(w, s); err != nil {
			return err
		}
	}

	return nil
}

func templateReader(conf *config.Config, path string, r io.Reader, size int64) (string, error) {
	v := visitor.NewTemplateComments(conf, path)

	decoder := yaml.NewDecoder(r)
	var buf strings.Builder
	if size != 0 {
		buf.Grow(int(size))
	}

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

	return buf.String(), nil
}
