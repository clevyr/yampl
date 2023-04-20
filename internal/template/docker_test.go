package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerRepo(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"postgres:alpine", args{"postgres:alpine"}, "postgres"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DockerRepo(tt.args.image)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDockerTag(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"postgres:alpine", args{"postgres:alpine"}, "alpine"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DockerTag(tt.args.image)
			assert.Equal(t, tt.want, got)
		})
	}
}
