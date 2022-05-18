package template

import "testing"

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
			if got := DockerRepo(tt.args.image); got != tt.want {
				t.Errorf("DockerRepo() = %v, want %v", got, tt.want)
			}
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
			if got := DockerTag(tt.args.image); got != tt.want {
				t.Errorf("DockerTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
