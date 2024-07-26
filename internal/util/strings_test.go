package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowerFirst(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"len 1", args{"A"}, "a"},
		{"multiple", args{"TestArg"}, "testArg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, LowerFirst(tt.args.s), "LowerFirst(%v)", tt.args.s)
		})
	}
}

func TestUpperFirst(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"len 1", args{"a"}, "A"},
		{"multiple", args{"testArg"}, "TestArg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UpperFirst(tt.args.s), "UpperFirst(%v)", tt.args.s)
		})
	}
}
