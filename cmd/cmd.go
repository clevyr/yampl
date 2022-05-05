package cmd

import (
	"errors"
	"fmt"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/pkg/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

var Command = &cobra.Command{
	Use:                   "yampl [-i] [-p prefix] [-v key=value ...] [file ...]",
	Short:                 "Inline YAML templating via line comments",
	DisableFlagsInUseLine: true,
	DisableAutoGenTag:     true,
	PreRunE:               preRun,
	RunE:                  run,
}

var conf config.Config

func init() {
	Command.Flags().StringToStringVarP(&conf.Values, "value", "v", map[string]string{}, "Define a template variable")
	Command.Flags().BoolVarP(&conf.Inline, "inline", "i", false, "Edit files in-place")
	Command.Flags().StringVarP(&conf.Prefix, "prefix", "p", "#yampl", "Template prefix. Must begin with '#'")
	Command.Flags().StringVar(&conf.LeftDelim, "left-delim", "{{", "Override the default left delimiter")
	Command.Flags().StringVar(&conf.RightDelim, "right-delim", "}}", "Override the default right delimiter")
}

func preRun(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	conf.Paths = args
	if !strings.HasPrefix(conf.Prefix, "#") {
		return errors.New("prefix must begin with '#'")
	}

	if conf.Inline && len(args) == 0 {
		return errors.New("no input files")
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		b, err := templateReader(os.Stdin)
		if err != nil {
			return err
		}

		fmt.Print(string(b))
	}

	for i, p := range args {
		if err := openAndTemplate(p); err != nil {
			return err
		}

		if !conf.Inline && i != len(args)-1 {
			fmt.Println("---")
		}
	}

	return nil
}

func openAndTemplate(p string) (err error) {
	var f *os.File
	if conf.Inline {
		stat, err := os.Stat(p)
		if err != nil {
			return err
		}

		f, err = os.OpenFile(p, os.O_RDWR, stat.Mode())
	} else {
		f, err = os.Open(p)
	}
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	b, err := templateReader(f)
	if err != nil {
		return err
	}

	if conf.Inline {
		if err := f.Truncate(int64(len(b))); err != nil {
			return err
		}

		if _, err := f.Write(b); err != nil {
			return err
		}
	} else {
		fmt.Print(string(b))
	}

	return f.Close()
}

func templateReader(r io.Reader) ([]byte, error) {
	t := template.LineComment{
		Config: conf,
	}

	if err := yaml.NewDecoder(r).Decode(&t); err != nil {
		return []byte{}, err
	}

	b, err := yaml.Marshal(t)
	if err != nil {
		return b, err
	}

	return b, nil
}
