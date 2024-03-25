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

func Test_preRun(t *testing.T) {
	t.Run("silent usage", func(t *testing.T) {
		cmd := NewCommand("", "")
		_ = preRun(cmd, []string{})
		assert.True(t, cmd.SilenceUsage)
	})

	t.Run("no error", func(t *testing.T) {
		err := preRun(NewCommand("", ""), []string{})
		require.NoError(t, err)
	})

	t.Run("invalid prefix", func(t *testing.T) {
		conf.Prefix = "tmpl"
		defer func() {
			conf.Prefix = "#yampl"
		}()

		if err := preRun(NewCommand("", ""), []string{}); !assert.NoError(t, err) {
			return
		}

		want := "#tmpl"
		assert.Equal(t, want, conf.Prefix)
	})

	t.Run("inplace no files", func(t *testing.T) {
		conf.Inplace = true
		defer func() {
			conf.Inplace = false
		}()

		err := preRun(NewCommand("", ""), []string{})
		require.Error(t, err)
	})

	t.Run("completion flag enabled", func(t *testing.T) {
		cmd := NewCommand("", "")
		if err := cmd.Flags().Set(CompletionFlag, "zsh"); !assert.NoError(t, err) {
			return
		}
		err := preRun(NewCommand("", ""), []string{})
		require.NoError(t, err)
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
	failConf := config.New()
	failConf.Fail = true

	stripConf := config.New()
	stripConf.Strip = true

	type args struct {
		conf config.Config
		r    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"empty", args{conf, strings.NewReader("")}, "", require.NoError},
		{"static", args{conf, strings.NewReader("a: a")}, "a: a\n", require.NoError},
		{"simple", args{conf, strings.NewReader("a: a #yampl b")}, "a: b #yampl b\n", require.NoError},
		{"dynamic", args{conf, strings.NewReader("a: a #yampl {{ .Value }}")}, "a: a #yampl {{ .Value }}\n", require.NoError},
		{"multi-doc", args{conf, strings.NewReader("a: a\n---\nb: b")}, "a: a\n---\nb: b\n", require.NoError},
		{"map flow to block", args{conf, strings.NewReader("a: {} #yampl:map {a: a}")}, "a: #yampl:map {a: a}\n  a: a\n", require.NoError},
		{"map block to flow", args{conf, strings.NewReader("a: #yampl:map {}\n  a: a")}, "a: {} #yampl:map {}\n", require.NoError},
		{"seq flow to block", args{conf, strings.NewReader("a: {} #yampl:seq [a]")}, "a: #yampl:seq [a]\n  - a\n", require.NoError},
		{"seq block to flow", args{conf, strings.NewReader("a: #yampl:seq []\n  - b")}, "a: [] #yampl:seq []\n", require.NoError},
		{"invalid yaml", args{conf, strings.NewReader("a:\n- b\n  c: c")}, "", require.Error},
		{"unset value allowed", args{conf, strings.NewReader("a: a #yampl {{ .b }}")}, "a: a #yampl {{ .b }}\n", require.NoError},
		{"unset value error", args{failConf, strings.NewReader("a: a #yampl {{ .z }}")}, "", require.Error},
		{"strip", args{stripConf, strings.NewReader("a: a #yampl b")}, "a: b\n", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := templateReader(tt.args.conf, tt.args.r)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_openAndTemplate(t *testing.T) {
	inplaceConf := config.New()
	inplaceConf.Inplace = true

	tempFileWith := func(contents string) (string, error) {
		f, err := os.CreateTemp("", "")
		if err != nil {
			return f.Name(), err
		}

		if _, err := f.WriteString(contents); err != nil {
			return f.Name(), err
		}

		if err := f.Close(); err != nil {
			return f.Name(), err
		}

		return f.Name(), nil
	}

	type args struct {
		conf     config.Config
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
			p, err := tempFileWith(tt.args.contents)
			defer func() {
				_ = os.RemoveAll(p)
			}()
			require.NoError(t, err)

			cmd := NewCommand("", "")
			var stdoutBuf strings.Builder
			cmd.SetOut(&stdoutBuf)

			err = openAndTemplate(cmd, tt.args.conf, p)
			tt.wantErr(t, err)

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
