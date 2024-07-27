package config

import (
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const EnvPrefix = "YAMPL_"

func (c *Config) Load(cmd *cobra.Command) error {
	IgnoredEnvs := []string{
		CompletionFlag,
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
		return errors.Join(errs...)
	}

	initLog(cmd)

	if !strings.HasPrefix(c.Prefix, "#") {
		c.Prefix = "#" + c.Prefix
	}

	c.Vars.Fill(c.valuesStringToString.Values())

	if f := cmd.Flags().Lookup(FailFlag); f.Changed {
		val, err := cmd.Flags().GetBool(FailFlag)
		if err != nil {
			return err
		}
		c.IgnoreUnsetErrors = !val
		c.IgnoreTemplateErrors = !val
	}

	return nil
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
