package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toYaml(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"map", args{map[string]any{"a": "a"}}, "a: a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toYaml(tt.args.v)
			assert.Equal(t, tt.want, got)
		})
	}
}
