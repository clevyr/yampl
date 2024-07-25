package cmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var version = "beta"

//go:embed description.md
var description string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "yampl [-i] [-p prefix] [-v key=value ...] [file ...]",
		Short:                 "Inline YAML templating via line-comments",
		Long:                  description,
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		ValidArgsFunction:     validArgs,
		Version:               buildVersion(version),
		RunE:                  run,
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

var ErrNoFiles = errors.New("no input files")

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

	if len(args) == 0 {
		if conf.Inplace || conf.Recursive {
			return ErrNoFiles
		}
		cmd.SilenceUsage = true

		if isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			return cmd.Help()
		}

		s, err := templateReader(conf, os.Stdin, log.Logger)
		if err != nil {
			return err
		}

		if err := colorize.WriteString(cmd.OutOrStdout(), s); err != nil {
			return err
		}
	}

	cmd.SilenceUsage = true

	for i, p := range args {
		if err := openAndTemplate(conf, cmd.OutOrStdout(), p); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}

		if !conf.Inplace && i != len(args)-1 {
			if _, err := io.WriteString(cmd.OutOrStdout(), "---\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

func openAndTemplate(conf *config.Config, w io.Writer, path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !util.IsYaml(path) {
				return err
			}

			if !conf.Inplace {
				if _, err := io.WriteString(w, "---\n"); err != nil {
					return err
				}
			}

			return openAndTemplateFile(conf, w, path)
		})
	}

	return openAndTemplateFile(conf, w, path)
}

func openAndTemplateFile(conf *config.Config, w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	log := log.With().Str("file", path).Logger()

	s, err := templateReader(conf, f, log)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
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
		if err := colorize.WriteString(w, s); err != nil {
			return err
		}
	}

	return nil
}

func templateReader(conf *config.Config, r io.Reader, log zerolog.Logger) (string, error) {
	v := visitor.NewTemplateComments(conf, log)

	decoder := yaml.NewDecoder(r)
	var buf strings.Builder

	for {
		var n yaml.Node

		if err := decoder.Decode(&n); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return buf.String(), err
		}

		if buf.Len() > 0 {
			buf.WriteString("---\n")
		}

		if err := v.Run(&n); err != nil {
			return buf.String(), err
		}

		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(conf.Indent)
		if err := encoder.Encode(&n); err != nil {
			_ = encoder.Close()
			return buf.String(), err
		}

		if err := encoder.Close(); err != nil {
			return buf.String(), err
		}
	}

	return buf.String(), nil
}
