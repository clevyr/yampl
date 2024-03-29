package cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/config/flags"
	"github.com/clevyr/yampl/internal/util"
	"github.com/clevyr/yampl/internal/visitor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

//go:embed description.md
var description string

func NewCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "yampl [-i] [-p prefix] [-v key=value ...] [file ...]",
		Short:                 "Inline YAML templating via line-comments",
		Long:                  description,
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		ValidArgsFunction:     validArgs,
		Version:               buildVersion(version, commit),
		PreRunE:               preRun,
		RunE:                  run,
	}

	registerCompletionFlag(cmd)
	registerLogFlags(cmd)
	registerValuesFlag(cmd)

	cmd.Flags().BoolVarP(&conf.Inplace, "inplace", "i", conf.Inplace, "Edit files in place")
	if err := cmd.RegisterFlagCompletionFunc("inplace", util.BoolCompletion); err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&conf.Prefix, "prefix", "p", conf.Prefix, "Template comments must begin with this prefix. The beginning '#' is implied.")
	if err := cmd.RegisterFlagCompletionFunc("prefix", cobra.NoFileCompletions); err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&conf.LeftDelim, "left-delim", conf.LeftDelim, "Override template left delimiter")
	if err := cmd.RegisterFlagCompletionFunc("left-delim", cobra.NoFileCompletions); err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&conf.RightDelim, "right-delim", conf.RightDelim, "Override template right delimiter")
	if err := cmd.RegisterFlagCompletionFunc("right-delim", cobra.NoFileCompletions); err != nil {
		panic(err)
	}

	cmd.Flags().IntVarP(&conf.Indent, "indent", "I", conf.Indent, "Override output indentation")
	if err := cmd.RegisterFlagCompletionFunc("indent", cobra.NoFileCompletions); err != nil {
		panic(err)
	}

	cmd.Flags().BoolVarP(&conf.Fail, "fail", "f", conf.Fail, `Exit with an error if a template variable is not set`)
	if err := cmd.RegisterFlagCompletionFunc("fail", util.BoolCompletion); err != nil {
		panic(err)
	}

	cmd.Flags().BoolVarP(&conf.Strip, "strip", "s", conf.Strip, "Strip template comments from output")
	if err := cmd.RegisterFlagCompletionFunc("strip", util.BoolCompletion); err != nil {
		panic(err)
	}

	cmd.InitDefaultVersionFlag()

	return cmd
}

//nolint:gochecknoglobals
var conf = config.New()

func validArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
}

var ErrNoFiles = errors.New("no input files")

func preRun(cmd *cobra.Command, args []string) error {
	completionFlag, err := cmd.Flags().GetString(CompletionFlag)
	if err != nil {
		panic(err)
	}
	if completionFlag != "" {
		return nil
	}

	initLog(cmd)

	cmd.SilenceUsage = true

	if !strings.HasPrefix(conf.Prefix, "#") {
		conf.Prefix = "#" + conf.Prefix
	}

	if conf.Inplace && len(args) == 0 {
		return ErrNoFiles
	}

	rawValues, err := cmd.Flags().GetStringToString(flags.ValueFlag)
	if err != nil {
		panic(err)
	}

	conf.Values.Fill(rawValues)

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	completionFlag, err := cmd.Flags().GetString(CompletionFlag)
	if err != nil {
		panic(err)
	}
	if completionFlag != "" {
		return completion(cmd, args)
	}

	if len(args) == 0 {
		s, err := templateReader(conf, os.Stdin)
		if err != nil {
			return err
		}

		//nolint:forbidigo
		fmt.Print(s)
	}

	for i, p := range args {
		if err := openAndTemplate(cmd, conf, p); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}

		if !conf.Inplace && i != len(args)-1 {
			//nolint:forbidigo
			fmt.Println("---")
		}
	}

	return nil
}

func openAndTemplate(cmd *cobra.Command, conf config.Config, p string) error {
	defer func(logger *log.Entry) {
		conf.Log = logger
	}(conf.Log)
	conf.Log = log.WithField("file", p)

	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	s, err := templateReader(conf, f)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	_ = f.Close()

	if conf.Inplace {
		temp, err := os.CreateTemp("", "yampl_*_"+filepath.Base(p))
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

		if err := os.Rename(temp.Name(), p); err != nil {
			log.Trace("failed to rename file, attempting to copy contents")

			in, err := os.Open(temp.Name())
			if err != nil {
				return err
			}
			defer func() {
				_ = in.Close()
			}()

			out, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC, stat.Mode())
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
		if _, err := fmt.Fprint(cmd.OutOrStdout(), s); err != nil {
			return err
		}
	}

	return nil
}

func templateReader(conf config.Config, r io.Reader) (string, error) {
	v := visitor.NewTemplateComments(conf)

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
