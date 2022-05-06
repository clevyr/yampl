package template

import "strings"

func DockerRepo(image string) string {
	split := strings.SplitN(image, ":", 2)
	return split[0]
}

func DockerTag(image string) string {
	split := strings.SplitN(image, ":", 2)
	return split[1]
}
