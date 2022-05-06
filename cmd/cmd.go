package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	"github.com/clevyr/go-yampl/internal/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

var Command = &cobra.Command{
	Use:                   "yampl [-i] [-p prefix] [-v key=value ...] [file ...]",
	Short:                 "Inline YAML templating via line-comments",
	DisableFlagsInUseLine: true,
	DisableAutoGenTag:     true,
	ValidArgsFunction:     validArgs,
	Version:               buildVersion(),
	PreRunE:               preRun,
	RunE:                  run,
}

var conf config.Config

func init() {
	Command.Flags().StringToStringVarP((*map[string]string)(&conf.Values), "value", "v", map[string]string{}, "Define a template variable")
	Command.Flags().BoolVarP(&conf.Inline, "inline", "i", false, "Edit files in-place")
	Command.Flags().StringVarP(&conf.Prefix, "prefix", "p", "#yampl", "Template prefix. Must begin with '#'")
	Command.Flags().StringVar(&conf.LeftDelim, "left-delim", "{{", "Override the default left delimiter")
	Command.Flags().StringVar(&conf.RightDelim, "right-delim", "}}", "Override the default right delimiter")
}

func validArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
}

func preRun(cmd *cobra.Command, args []string) error {
	if completionFlag != "" {
		return nil
	}

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
	if completionFlag != "" {
		return completion(cmd, args)
	}

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
		if err != nil {
			return err
		}
	} else {
		f, err = os.Open(p)
		if err != nil {
			return err
		}
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

		if _, err := f.Seek(0, io.SeekStart); err != nil {
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
	decoder := yaml.NewDecoder(r)
	var buf bytes.Buffer

	for {
		var n yaml.Node

		if err := decoder.Decode(&n); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return buf.Bytes(), err
		}

		if buf.Len() > 0 {
			buf.Write([]byte("---\n"))
		}

		if err := node.Visit(conf, template.LineComment, &n); err != nil {
			return buf.Bytes(), err
		}

		b, err := yaml.Marshal(&n)
		if err != nil {
			return buf.Bytes(), err
		}

		buf.Write(b)
	}

	return buf.Bytes(), nil
}
