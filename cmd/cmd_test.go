package cmd

import (
	"github.com/spf13/cobra"
	"reflect"
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
		if err := preRun(&cobra.Command{}, []string{}); err != nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("invalid prefix", func(t *testing.T) {
		conf.Prefix = "a"
		defer func() {
			conf.Prefix = "#yampl"
		}()

		if err := preRun(&cobra.Command{}, []string{}); err == nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("inline no files", func(t *testing.T) {
		conf.Inline = true
		defer func() {
			conf.Inline = false
		}()

		if err := preRun(&cobra.Command{}, []string{}); err == nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("completion flag enabled", func(t *testing.T) {
		completionFlag = "zsh"
		defer func() {
			completionFlag = ""
		}()
		if err := preRun(&cobra.Command{}, []string{}); err != nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})
}

func Test_validArgs(t *testing.T) {
	type args struct {
		cmd        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 cobra.ShellCompDirective
	}{
		{"default", args{}, []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validArgs(tt.args.cmd, tt.args.args, tt.args.toComplete)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validArgs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("validArgs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
