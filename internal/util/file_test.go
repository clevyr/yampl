package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsYaml(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"yaml", args{"test.yaml"}, true},
		{"yml", args{"test.yml"}, true},
		{"txt", args{"test.txt"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsYaml(tt.args.path))
		})
	}
}
