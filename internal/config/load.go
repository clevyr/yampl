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
		if !f.Changed && !slices.Contains(IgnoredEnvs, f.Name) {
			if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
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

	rawValues, err := cmd.Flags().GetStringToString(ValueFlag)
	if err != nil {
		return err
	}
	c.Values.Fill(rawValues)

	return nil
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
