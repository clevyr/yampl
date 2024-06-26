package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		{"newline with equal", args{[]string{"yampl", "--value=a=a\nb=b"}}, []string{"yampl", "--value=a=a", "--value=b=b"}},
		{"newline with space", args{[]string{"yampl", "--value", "a=a\nb=b"}}, []string{"yampl", "--value=a=a", "--value=b=b"}},
		{"newline in file", args{[]string{"yampl", "test\nfile.yaml"}}, []string{"yampl", "test\nfile.yaml"}},
		{"newline after end of options", args{[]string{"yampl", "-v=a=a", "---", "-v\nfile.yaml"}}, []string{"yampl", "-v=a=a", "---", "-v\nfile.yaml"}},
		{"trim newline", args{[]string{"yampl", "-v=\na=a\n"}}, []string{"yampl", "-v=a=a"}},
		{"collapse newlines", args{[]string{"yampl", "-v=a=a\n\nb=b"}}, []string{"yampl", "-v=a=a", "-v=b=b"}},
		{"trim spaces", args{[]string{"yampl", "-v=a=a\n  b=b"}}, []string{"yampl", "-v=a=a", "-v=b=b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FixStringToStringNewlines(tt.args.s)
			assert.Equal(t, tt.want, got)
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
			got := hasValueFlag(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
