package cmd

import (
	"io"
	"testing"

	"github.com/clevyr/yampl/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_completion(t *testing.T) {
	r, w := io.Pipe()
	_ = r.Close()

	type args struct {
		cmd   *cobra.Command
		shell string
	}
	tests := []struct {
		name    string
		w       io.Writer
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{"bash", io.Discard, args{NewCommand(), "bash"}, require.NoError},
		{"bash error", w, args{NewCommand(), "bash"}, require.Error},
		{"zsh", io.Discard, args{NewCommand(), "zsh"}, require.NoError},
		{"zsh error", w, args{NewCommand(), "zsh"}, require.Error},
		{"fish", io.Discard, args{NewCommand(), "fish"}, require.NoError},
		{"fish error", w, args{NewCommand(), "fish"}, require.Error},
		{"powershell", io.Discard, args{NewCommand(), "powershell"}, require.NoError},
		{"powershell error", w, args{NewCommand(), "powershell"}, require.Error},
		{"other", io.Discard, args{NewCommand(), "other"}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.cmd.SetOut(tt.w)

			if err := tt.args.cmd.Flags().Set(config.CompletionFlag, tt.args.shell); !assert.NoError(t, err) {
				return
			}
			err := completion(tt.args.cmd, tt.args.shell)
			tt.wantErr(t, err)
		})
	}
}
