package config

import (
	"errors"
	"os"
	"slices"
	"strings"

	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const EnvPrefix = "YAMPL_"

var ErrCmdMissingConfig = errors.New("config missing from command context")

func Load(cmd *cobra.Command) (*Config, error) {
	conf, ok := FromContext(cmd.Context())
	if !ok {
		return nil, ErrCmdMissingConfig
	}

	IgnoredEnvs := []string{
		cobrax.FlagCompletion,
	}

	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if !slices.Contains(IgnoredEnvs, f.Name) {
				if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
					if err := f.Value.Set(val); err != nil {
						errs = append(errs, err)
					}
				}
			}
		}
	})
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	conf.InitLog(cmd.ErrOrStderr())

	if !strings.HasPrefix(conf.Prefix, "#") {
		conf.Prefix = "#" + conf.Prefix
	}

	conf.Vars.Fill(conf.valuesStringToString.Value())

	if f := cmd.Flags().Lookup(FailFlag); f.Changed {
		val, err := cmd.Flags().GetBool(FailFlag)
		if err != nil {
			return nil, err
		}
		conf.IgnoreUnsetErrors = !val
		conf.IgnoreTemplateErrors = !val
	}

	return conf, nil
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
