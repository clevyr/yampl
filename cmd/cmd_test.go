package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/clevyr/yampl/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_preRun(t *testing.T) {
	t.Run("silent usage", func(t *testing.T) {
		cmd := NewCommand("", "")
		_ = preRun(cmd, []string{})
		assert.True(t, cmd.SilenceUsage)
	})

	t.Run("no error", func(t *testing.T) {
		err := preRun(NewCommand("", ""), []string{})
		assert.NoError(t, err)
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
		assert.Error(t, err)
	})

	t.Run("completion flag enabled", func(t *testing.T) {
		cmd := NewCommand("", "")
		if err := cmd.Flags().Set(CompletionFlag, "zsh"); !assert.NoError(t, err) {
			return
		}
		err := preRun(NewCommand("", ""), []string{})
		assert.NoError(t, err)
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
		wantErr bool
	}{
		{"empty", args{conf, strings.NewReader("")}, "", false},
		{"static", args{conf, strings.NewReader("a: a")}, "a: a\n", false},
		{"simple", args{conf, strings.NewReader("a: a #yampl b")}, "a: b #yampl b\n", false},
		{"dynamic", args{conf, strings.NewReader("a: a #yampl {{ .Value }}")}, "a: a #yampl {{ .Value }}\n", false},
		{"multi-doc", args{conf, strings.NewReader("a: a\n---\nb: b")}, "a: a\n---\nb: b\n", false},
		{"map flow to block", args{conf, strings.NewReader("a: {} #yampl:map {a: a}")}, "a: #yampl:map {a: a}\n  a: a\n", false},
		{"map block to flow", args{conf, strings.NewReader("a: #yampl:map {}\n  a: a")}, "a: {} #yampl:map {}\n", false},
		{"seq flow to block", args{conf, strings.NewReader("a: {} #yampl:seq [a]")}, "a: #yampl:seq [a]\n  - a\n", false},
		{"seq block to flow", args{conf, strings.NewReader("a: #yampl:seq []\n  - b")}, "a: [] #yampl:seq []\n", false},
		{"invalid yaml", args{conf, strings.NewReader("a:\n- b\n  c: c")}, "", true},
		{"unset value allowed", args{conf, strings.NewReader("a: a #yampl {{ .b }}")}, "a: a #yampl {{ .b }}\n", false},
		{"unset value error", args{failConf, strings.NewReader("a: a #yampl {{ .z }}")}, "", true},
		{"strip", args{stripConf, strings.NewReader("a: a #yampl b")}, "a: b\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := templateReader(tt.args.conf, tt.args.r)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_openAndTemplate(t *testing.T) {
	inplaceConf := config.New()
	inplaceConf.Inplace = true

	tempFileWith := func(contents string) (*os.File, func(), error) {
		f, err := os.CreateTemp("", "")
		if err != nil {
			return nil, func() {}, err
		}

		if _, err := f.WriteString(contents); err != nil {
			return nil, func() {}, err
		}

		return f, func() {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}, nil
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
		wantErr    bool
	}{
		{"simple", args{config.New(), "a: a"}, "a: a\n", true, false},
		{"template", args{config.New(), "a: a #yampl b"}, "a: b #yampl b\n", true, false},
		{"inplace", args{inplaceConf, "a: a #yampl b"}, "a: b #yampl b\n", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, cleanup, err := tempFileWith(tt.args.contents)
			if !assert.NoError(t, err) {
				return
			}
			defer cleanup()

			cmd := NewCommand("", "")
			var stdoutBuf strings.Builder
			cmd.SetOut(&stdoutBuf)

			if err := openAndTemplate(cmd, tt.args.conf, f.Name()); !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}

			if _, err := f.Seek(0, io.SeekStart); !assert.NoError(t, err) {
				return
			}

			var buf strings.Builder
			if _, err := io.Copy(&buf, f); !assert.NoError(t, err) {
				return
			}

			if tt.wantStdout {
				assert.Equal(t, tt.want, stdoutBuf.String())
				assert.Equal(t, tt.args.contents, buf.String())
			} else {
				assert.Empty(t, stdoutBuf.String())
				assert.Equal(t, tt.want, buf.String())
			}
		})
	}
}
