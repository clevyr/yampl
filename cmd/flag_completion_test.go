package cmd

import (
	"github.com/spf13/cobra"
	"reflect"
	"testing"
)

func Test_completion(t *testing.T) {
	type args struct {
		cmd   *cobra.Command
		args  []string
		shell string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"bash", args{&cobra.Command{}, []string{}, "bash"}, false},
		{"zsh", args{&cobra.Command{}, []string{}, "zsh"}, false},
		{"fish", args{&cobra.Command{}, []string{}, "fish"}, false},
		{"powershell", args{&cobra.Command{}, []string{}, "powershell"}, false},
		{"other", args{&cobra.Command{}, []string{}, "other"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
