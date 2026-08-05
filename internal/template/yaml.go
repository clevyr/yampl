package template

import (
	"bytes"

	"github.com/goccy/go-yaml"
)

func toYaml(v any) (string, error) {
	data, err := yaml.MarshalWithOptions(v, yaml.IndentSequence(true))
	if err != nil {
		return "", err
	}

	// IndentSequence also indents the root sequence; shift it back so
	// top-level lists start at column zero
	if bytes.HasPrefix(data, []byte("  - ")) {
		lines := bytes.Split(data, []byte("\n"))
		for i, line := range lines {
			lines[i] = bytes.TrimPrefix(line, []byte("  "))
		}
		data = bytes.Join(lines, []byte("\n"))
	}

	return string(bytes.TrimSuffix(data, []byte("\n"))), nil
}
