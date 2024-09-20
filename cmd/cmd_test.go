package cmd

import (
	"io"
	"os"
	"strings"
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

	t.Run("completion flag enabled", func(t *testing.T) {
		cmd := New()
		if err := cmd.Flags().Set(config.CompletionFlag, "zsh"); !assert.NoError(t, err) {
			return
		}
		require.NoError(t, run(cmd, []string{}))
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

func Test_templateReader(t *testing.T) {
	ignoreTemplateConf := config.New()
	ignoreTemplateConf.IgnoreTemplateErrors = true

	failUnsetConf := config.New()
	failUnsetConf.IgnoreUnsetErrors = false

	stripConf := config.New()
	stripConf.Strip = true

	type args struct {
		conf *config.Config
		r    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"empty", args{config.New(), strings.NewReader("")}, "", require.NoError},
		{"static", args{config.New(), strings.NewReader("a: a")}, "a: a\n", require.NoError},
		{"simple", args{config.New(), strings.NewReader("a: a #yampl b")}, "a: b #yampl b\n", require.NoError},
		{"dynamic (deprecated)", args{config.New(), strings.NewReader("a: a #yampl {{ upper .Value }}")}, "a: A #yampl {{ upper .Value }}\n", require.NoError},
		{"dynamic", args{config.New(), strings.NewReader("a: a #yampl {{ upper current }}")}, "a: A #yampl {{ upper current }}\n", require.NoError},
		{"multi-doc", args{config.New(), strings.NewReader("a: a\n---\nb: b")}, "a: a\n---\nb: b\n", require.NoError},
		{"map flow to block", args{config.New(), strings.NewReader("a: {} #yampl:map {a: a}")}, "a: #yampl:map {a: a}\n  a: a\n", require.NoError},
		{"map block to flow", args{config.New(), strings.NewReader("a: #yampl:map {}\n  a: a")}, "a: {} #yampl:map {}\n", require.NoError},
		{"seq flow to block", args{config.New(), strings.NewReader("a: {} #yampl:seq [a]")}, "a: #yampl:seq [a]\n  - a\n", require.NoError},
		{"seq block to flow", args{config.New(), strings.NewReader("a: #yampl:seq []\n  - b")}, "a: [] #yampl:seq []\n", require.NoError},
		{"invalid yaml", args{config.New(), strings.NewReader("a:\n- b\n  c: c")}, "", require.Error},
		{"invalid template", args{config.New(), strings.NewReader("a: a #yampl {{ current")}, "", require.Error},
		{"invalid template ignored", args{ignoreTemplateConf, strings.NewReader("a: a #yampl {{ current")}, "a: a #yampl {{ current\n", require.NoError},
		{"unset value ignored", args{config.New(), strings.NewReader("a: a #yampl {{ .b }}")}, "a: a #yampl {{ .b }}\n", require.NoError},
		{"unset value error", args{failUnsetConf, strings.NewReader("a: a #yampl {{ .z }}")}, "", require.Error},
		{"strip", args{stripConf, strings.NewReader("a: a #yampl b")}, "a: b\n", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var size int64
			if r, ok := tt.args.r.(*strings.Reader); ok {
				size = r.Size()
			}
			got, err := templateReader(tt.args.conf, "", tt.args.r, size)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_openAndTemplateFile(t *testing.T) {
	inplaceConf := config.New()
	inplaceConf.Inplace = true

	tempFile := func(t *testing.T, contents string) string {
		f, err := os.CreateTemp("", "")
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = f.Close()
			_ = os.Remove(f.Name())
		})

		_, err = f.WriteString(contents)
		require.NoError(t, err)
		require.NoError(t, f.Close())
		return f.Name()
	}

	type args struct {
		conf     *config.Config
		contents string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStdout bool
		wantErr    require.ErrorAssertionFunc
	}{
		{"simple", args{config.New(), "a: a"}, "a: a\n", true, require.NoError},
		{"template", args{config.New(), "a: a #yampl b"}, "a: b #yampl b\n", true, require.NoError},
		{"inplace", args{inplaceConf, "a: a #yampl b"}, "a: b #yampl b\n", false, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tempFile(t, tt.args.contents)

			var stdoutBuf strings.Builder
			tt.wantErr(t, openAndTemplateFile(tt.args.conf, &stdoutBuf, p, p, false))

			fileContents, err := os.ReadFile(p)
			require.NoError(t, err)

			if tt.wantStdout {
				assert.Equal(t, tt.want, stdoutBuf.String())
				assert.EqualValues(t, tt.args.contents, fileContents)
			} else {
				assert.Empty(t, stdoutBuf.String())
				assert.EqualValues(t, tt.want, fileContents)
			}
		})
	}
}
