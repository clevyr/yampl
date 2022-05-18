package cmd

import (
	"github.com/clevyr/go-yampl/internal/visitor"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/spf13/cobra"
	"os"
)

var rawValues map[string]string

func init() {
	Command.Flags().StringToStringVarP(&rawValues, "value", "v", map[string]string{}, "Define a template variable. Can be used more than once.")
	err := Command.RegisterFlagCompletionFunc("value", valueCompletion)
	if err != nil {
		panic(err)
	}
}

func valueCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	v := visitor.NewFindArgs(conf)

	for _, path := range args {
		func() {
			b, err := os.ReadFile(path)
			if err != nil {
				return
			}

			file, err := parser.ParseBytes(b, parser.ParseComments)
			if err != nil {
				return
			}

			for _, doc := range file.Docs {
				if ast.Walk(&v, doc.Body); v.Error() != nil {
					return
				}
			}
		}()
	}

	return v.Values(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
