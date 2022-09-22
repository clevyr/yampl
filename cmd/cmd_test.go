package cmd

import (
	"github.com/clevyr/yampl/internal/config"
	"github.com/spf13/cobra"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_preRun(t *testing.T) {
	t.Run("silent usage", func(t *testing.T) {
		cmd := Command
		_ = preRun(cmd, []string{})
		if !cmd.SilenceUsage {
			t.Errorf("preRun() Command.SilenceUsage got = %v, want %v", cmd.SilenceUsage, false)
		}
	})

	t.Run("no error", func(t *testing.T) {
		if err := preRun(Command, []string{}); err != nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("invalid prefix", func(t *testing.T) {
		conf.Prefix = "tmpl"
		defer func() {
			conf.Prefix = "#yampl"
		}()

		if err := preRun(&cobra.Command{}, []string{}); err != nil {
			t.Errorf("preRun() error = %v, wantErr %v", err, false)
		}

		want := "#tmpl"
		if conf.Prefix != want {
			t.Errorf("preRun() prefix = %s, want %s", conf.Prefix, want)
		}
	})

	t.Run("inplace no files", func(t *testing.T) {
		conf.Inplace = true
		defer func() {
			conf.Inplace = false
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
			if (err != nil) != tt.wantErr {
				t.Errorf("templateReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("templateReader() got = %v, want %v", string(got), string(tt.want))
			}
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
		{"simple", args{config.New(), "a: a"}, "a: a", true, false},
		{"template", args{config.New(), "a: a #yampl b"}, "a: a #yampl b", true, false},
		{"inplace", args{inplaceConf, "a: a #yampl b"}, "a: b #yampl b\n", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Error(err)
				return
			}

			var stdoutBuf strings.Builder
			ch := make(chan struct{})
			go func() {
				_, _ = io.Copy(&stdoutBuf, r)
				ch <- struct{}{}
			}()
			defer func(w *os.File) {
				_ = w.Close()
			}(w)

			defer func(stdout *os.File) {
				os.Stdout = stdout
			}(os.Stdout)
			os.Stdout = w

			f, cleanup, err := tempFileWith(tt.args.contents)
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup()

			if err := openAndTemplate(tt.args.conf, f.Name()); (err != nil) != tt.wantErr {
				t.Errorf("openAndTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if _, err := f.Seek(0, io.SeekStart); err != nil {
				t.Error(err)
				return
			}

			var buf strings.Builder
			if _, err := io.Copy(&buf, f); err != nil {
				t.Error(err)
				return
			}

			_ = w.Close()
			<-ch

			if (stdoutBuf.Len() != 0) != tt.wantStdout {
				t.Errorf("openAndTemplate() got stdout len = %v, want stdout %v", stdoutBuf.Len(), tt.wantStdout)
				return
			}

			var got string
			if tt.wantStdout {
				got = stdoutBuf.String()
			} else {
				got = buf.String()
			}
			if buf.String() != tt.want {
				t.Errorf("openAndTemplate() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
