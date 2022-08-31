package cmd

import (
	"github.com/clevyr/go-yampl/internal/node"
	"github.com/clevyr/go-yampl/internal/visitor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var rawValues map[string]string

func init() {
	Command.Flags().StringToStringVarP(&rawValues, "value", "v", rawValues, "Define a template variable. Can be used more than once.")
	err := Command.RegisterFlagCompletionFunc("value", valueCompletion)
	if err != nil {
		panic(err)
	}
}

func valueCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if !strings.HasPrefix(conf.Prefix, "#") {
		conf.Prefix = "#" + conf.Prefix
	}
	conf.Prefix += " "

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
				continue
			}

			if err := node.Visit(v.Visit, &n); err != nil {
				continue
			}
		}
	}

	return v.Values(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
