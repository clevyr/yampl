package util

import (
	"reflect"
	"testing"
)

func TestFixStringToStringNewlines(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"no value flag", args{[]string{"yampl"}}, []string{"yampl"}},
		{"no newline", args{[]string{"yampl", "--value=a=a"}}, []string{"yampl", "--value=a=a"}},
		{"newline with equal", args{[]string{"yampl", "--value=a=a\nb=b"}}, []string{"yampl", "--value=a=a,b=b"}},
		{"newline with space", args{[]string{"yampl", "--value", "a=a\nb=b"}}, []string{"yampl", "--value", "a=a,b=b"}},
		{"newline in file", args{[]string{"yampl", "test\nfile.yaml"}}, []string{"yampl", "test\nfile.yaml"}},
		{"newline after end of options", args{[]string{"yampl", "-v=a=a", "---", "-v\nfile.yaml"}}, []string{"yampl", "-v=a=a", "---", "-v\nfile.yaml"}},
		{"trim newline", args{[]string{"yampl", "-v=\na=a\n"}}, []string{"yampl", "-v=a=a"}},
		{"collapse newlines", args{[]string{"yampl", "-v=a=a\n\nb=b"}}, []string{"yampl", "-v=a=a,b=b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FixStringToStringNewlines(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FixStringToStringNewlines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasValueFlag(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"no flag", args{"yampl"}, false},
		{"normal", args{"--value"}, true},
		{"normal with value", args{"--value=test"}, true},
		{"shorthand", args{"-v"}, true},
		{"shorthand with value", args{"-v=test"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasValueFlag(tt.args.s); got != tt.want {
				t.Errorf("hasValueFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
