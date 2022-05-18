package template

import (
	"strings"
)

func DockerRepo(image string) string {
	repo, _, _ := strings.Cut(image, ":")
	return repo
}

func DockerTag(image string) string {
	_, tag, _ := strings.Cut(image, ":")
	return tag
}
