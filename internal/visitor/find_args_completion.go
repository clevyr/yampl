package visitor

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func RegisterCompletion(cmd *cobra.Command, conf *config.Config) {
	if err := cmd.RegisterFlagCompletionFunc(config.VarFlag, valueCompletion(conf)); err != nil {
		panic(err)
	}
}

func valueCompletion(conf *config.Config) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		_ = conf.Load(cmd)

		if !strings.HasPrefix(conf.Prefix, "#") {
			conf.Prefix = "#" + conf.Prefix
		}

		v := NewFindArgs(conf)
		for _, path := range args {
			if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() || !util.IsYaml(path) {
					return err
				}

				return valueCompletionFile(path, v)
			}); err != nil {
				continue
			}
		}

		return v.Values(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}

func valueCompletionFile(path string, v *FindArgs) error {
	v.path = path

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)

	for {
		var n yaml.Node

		if err := decoder.Decode(&n); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if err := v.Run(&n); err != nil {
			return err
		}
	}
}
