package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/clevyr/yampl/internal/visitor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rawValues map[string]string

const (
	ValueFlag      = "value"
	ValueFlagShort = "v"
)

func init() {
	Command.Flags().StringToStringVarP(&rawValues, ValueFlag, ValueFlagShort, rawValues, "Define a template variable. Can be used more than once.")
	err := Command.RegisterFlagCompletionFunc("value", valueCompletion)
	if err != nil {
		panic(err)
	}
}

func valueCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if !strings.HasPrefix(conf.Prefix, "#") {
		conf.Prefix = "#" + conf.Prefix
	}

	v := visitor.NewFindArgs(conf)

	for _, path := range args {
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		decoder := yaml.NewDecoder(f)

		for {
			var n yaml.Node

			if err := decoder.Decode(&n); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				continue
			}

			if err := v.Run(&n); err != nil {
				continue
			}
		}
	}

	return v.Values(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
