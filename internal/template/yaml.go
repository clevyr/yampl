package template

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func toYaml(v any) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
