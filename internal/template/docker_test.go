package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerRepo(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"postgres:alpine", args{"postgres:alpine"}, "postgres", require.NoError},
		{"postgres", args{"postgres"}, "postgres", require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DockerRepo(tt.args.image)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDockerTag(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"postgres:alpine", args{"postgres:alpine"}, "alpine", require.NoError},
		{"postgres", args{"postgres"}, "", require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DockerTag(tt.args.image)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
