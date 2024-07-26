package template

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func toYaml(v any) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSuffix(data, []byte("\n"))), nil
}
