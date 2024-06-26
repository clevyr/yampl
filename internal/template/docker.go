package template

import (
	"errors"
	"strings"
)

var (
	ErrRepoNotFound = errors.New("docker repo not found")
	ErrTagNotFound  = errors.New("docker tag not found")
)

func DockerRepo(image string) (string, error) {
	repo, _, found := strings.Cut(image, ":")
	if !found {
		return repo, ErrRepoNotFound
	}
	return repo, nil
}

func DockerTag(image string) (string, error) {
	_, tag, found := strings.Cut(image, ":")
	if !found {
		return tag, ErrTagNotFound
	}
	return tag, nil
}
