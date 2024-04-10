package cmd

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_completion(t *testing.T) {
	r, w := io.Pipe()
	_ = r.Close()

	type args struct {
		cmd   *cobra.Command
		args  []string
		shell string
	}
	tests := []struct {
		name    string
		w       io.Writer
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{"bash", io.Discard, args{NewCommand(), []string{}, "bash"}, require.NoError},
		{"bash error", w, args{NewCommand(), []string{}, "bash"}, require.Error},
		{"zsh", io.Discard, args{NewCommand(), []string{}, "zsh"}, require.NoError},
		{"zsh error", w, args{NewCommand(), []string{}, "zsh"}, require.Error},
		{"fish", io.Discard, args{NewCommand(), []string{}, "fish"}, require.NoError},
		{"fish error", w, args{NewCommand(), []string{}, "fish"}, require.Error},
		{"powershell", io.Discard, args{NewCommand(), []string{}, "powershell"}, require.NoError},
		{"powershell error", w, args{NewCommand(), []string{}, "powershell"}, require.Error},
		{"other", io.Discard, args{NewCommand(), []string{}, "other"}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.cmd.SetOut(tt.w)

			if err := tt.args.cmd.Flags().Set(CompletionFlag, tt.args.shell); !assert.NoError(t, err) {
				return
			}
			err := completion(tt.args.cmd, tt.args.args)
			tt.wantErr(t, err)
		})
	}
}

func Test_completionCompletion(t *testing.T) {
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
		{"default", args{}, []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := completionCompletion(tt.args.cmd, tt.args.args, tt.args.toComplete)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
