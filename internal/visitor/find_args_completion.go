package visitor

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/parser"
	"github.com/clevyr/yampl/internal/util"
	"github.com/goccy/go-yaml/ast"
	"github.com/spf13/cobra"
)

func RegisterCompletion(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc(config.VarFlag, valueCompletion); err != nil {
		panic(err)
	}
}

func valueCompletion(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	conf, err := config.Load(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	v := NewFindArgs(conf)
	for _, path := range args {
		if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !util.IsYaml(path) {
				return err
			}

			return valueCompletionFile(path, v)
		}); err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
	}

	return v.Values(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func valueCompletionFile(path string, v *FindArgs) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	file, err := parser.ParseReader(f)
	if err != nil {
		return err
	}

	for _, doc := range file.Docs {
		if doc.Body == nil {
			continue
		}

		ast.Walk(v, doc.Body)
	}

	return nil
}
