package cmd

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	"github.com/clevyr/go-yampl/internal/template"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

//go:embed description.md
var description string

var Command = &cobra.Command{
	Use:                   "yampl [-i] [-p prefix] [-v key=value ...] [file ...]",
	Short:                 "Inline YAML templating via line-comments",
	Long:                  description,
	DisableFlagsInUseLine: true,
	DisableAutoGenTag:     true,
	ValidArgsFunction:     validArgs,
	Version:               buildVersion(),
	PreRunE:               preRun,
	RunE:                  run,
}

var conf = config.New()

func init() {
	Command.Flags().BoolVarP(&conf.Inplace, "inplace", "i", conf.Inplace, "Update files inplace")
	Command.Flags().StringVarP(&conf.Prefix, "prefix", "p", conf.Prefix, "Line-comments are ignored unless this prefix is found. Prefix must begin with '#'")
	Command.Flags().StringVar(&conf.LeftDelim, "left-delim", conf.LeftDelim, "Override the left delimiter")
	Command.Flags().StringVar(&conf.RightDelim, "right-delim", conf.RightDelim, "Override the right delimiter")
	Command.Flags().IntVarP(&conf.Indent, "indent", "I", conf.Indent, "Override output indentation")
	Command.Flags().BoolVarP(&conf.Strict, "strict", "s", conf.Strict, `Trigger an error if a template variable is missing`)
}

func validArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
}

func preRun(cmd *cobra.Command, args []string) error {
	if completionFlag != "" {
		return nil
	}

	cmd.SilenceUsage = true

	if !strings.HasPrefix(conf.Prefix, "#") {
		return errors.New("prefix must begin with '#'")
	}

	if conf.Inplace && len(args) == 0 {
		return errors.New("no input files")
	}

	conf.Values.Fill(rawValues)

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	if completionFlag != "" {
		return completion(cmd, args)
	}

	if len(args) == 0 {
		b, err := templateReader(conf, os.Stdin)
		if err != nil {
			return err
		}

		fmt.Print(string(b))
	}

	for i, p := range args {
		if err := openAndTemplate(conf, p); err != nil {
			return err
		}

		if !conf.Inplace && i != len(args)-1 {
			fmt.Println("---")
		}
	}

	return nil
}

func openAndTemplate(conf config.Config, p string) (err error) {
	defer func(logger *log.Entry) {
		conf.Log = logger
	}(conf.Log)
	conf.Log = log.WithField("file", p)

	var f *os.File
	if conf.Inplace {
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

	b, err := templateReader(conf, f)
	if err != nil {
		return err
	}

	if conf.Inplace {
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

func templateReader(conf config.Config, r io.Reader) ([]byte, error) {
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

		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(conf.Indent)
		if err := encoder.Encode(&n); err != nil {
			_ = encoder.Close()
			return buf.Bytes(), err
		}

		if err := encoder.Close(); err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}
