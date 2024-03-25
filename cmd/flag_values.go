package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/clevyr/yampl/internal/config/flags"
	"github.com/clevyr/yampl/internal/visitor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func registerValuesFlag(cmd *cobra.Command) {
	cmd.Flags().StringToStringP(flags.ValueFlag, flags.ValueFlagShort, map[string]string{}, "Define a template variable. Can be used more than once.")
	err := cmd.RegisterFlagCompletionFunc("value", valueCompletion)
	if err != nil {
		panic(err)
	}
}

func valueCompletion(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
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
