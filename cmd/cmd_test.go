package cmd

import (
	"github.com/spf13/cobra"
	"testing"
)

func Test_preRun(t *testing.T) {
	t.Run("silent usage", func(t *testing.T) {
		var cmd cobra.Command
		_ = preRun(&cmd, []string{})
		if !cmd.SilenceUsage {
			t.Errorf("preRun() Command.SilenceUsage got = %v, want %v", cmd.SilenceUsage, false)
		}
	})

	t.Run("no error", func(t *testing.T) {
		var cmd cobra.Command
		if err := preRun(&cmd, []string{}); err != nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("invalid prefix", func(t *testing.T) {
		var cmd cobra.Command
		conf.Prefix = "a"
		defer func() {
			conf.Prefix = "#yampl"
		}()

		if err := preRun(&cmd, []string{}); err == nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("inline no files", func(t *testing.T) {
		var cmd cobra.Command
		conf.Inline = true
		defer func() {
			conf.Inline = false
		}()

		if err := preRun(&cmd, []string{}); err == nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})
}
