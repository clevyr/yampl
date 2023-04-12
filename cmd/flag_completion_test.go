package cmd

import (
	"io"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
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
		wantErr bool
	}{
		{"bash", io.Discard, args{&cobra.Command{}, []string{}, "bash"}, false},
		{"bash error", w, args{&cobra.Command{}, []string{}, "bash"}, true},
		{"zsh", io.Discard, args{&cobra.Command{}, []string{}, "zsh"}, false},
		{"zsh error", w, args{&cobra.Command{}, []string{}, "zsh"}, true},
		{"fish", io.Discard, args{&cobra.Command{}, []string{}, "fish"}, false},
		{"fish error", w, args{&cobra.Command{}, []string{}, "fish"}, true},
		{"powershell", io.Discard, args{&cobra.Command{}, []string{}, "powershell"}, false},
		{"powershell error", w, args{&cobra.Command{}, []string{}, "powershell"}, true},
		{"other", io.Discard, args{&cobra.Command{}, []string{}, "other"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(prev io.Writer) {
				completionWriter = prev
			}(completionWriter)
			completionWriter = tt.w

			completionFlag = tt.args.shell
			if err := completion(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("completion() error = %v, wantErr %v", err, tt.wantErr)
			}
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("completionCompletion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("completionCompletion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
