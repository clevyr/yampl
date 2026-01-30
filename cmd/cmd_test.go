package cmd

import (
	"testing"

	"github.com/clevyr/yampl/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_run(t *testing.T) {
	t.Run("silent usage", func(t *testing.T) {
		cmd := New()
		_ = run(cmd, []string{})
		assert.True(t, cmd.SilenceUsage)
	})

	t.Run("no error", func(t *testing.T) {
		require.NoError(t, run(New(), []string{}))
	})

	t.Run("invalid prefix", func(t *testing.T) {
		cmd := New()
		conf, ok := config.FromContext(cmd.Context())
		require.True(t, ok)
		conf.Prefix = "tmpl"
		require.NoError(t, run(cmd, []string{}))
		want := "#tmpl"
		assert.Equal(t, want, conf.Prefix)
	})

	t.Run("inplace no files", func(t *testing.T) {
		cmd := New()
		conf, ok := config.FromContext(cmd.Context())
		require.True(t, ok)
		conf.Inplace = true
		require.Error(t, run(cmd, []string{}))
	})

	t.Run("has config", func(t *testing.T) {
		cmd := New()
		conf, ok := config.FromContext(cmd.Context())
		assert.True(t, ok)
		assert.NotNil(t, conf)
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
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
