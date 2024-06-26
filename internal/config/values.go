package config

import (
	"strings"
)

type Values map[string]any

func (values Values) Fill(rawValues map[string]string) {
	for k, v := range rawValues {
		split := strings.Split(k, ".")
		if len(split) == 1 {
			values[k] = v
		} else {
			values.SetNested(v, split...)
		}
	}
}

func (values Values) SetNested(v any, k ...string) {
	setNested(values, v, k...)
}

func setNested(input Values, value any, keys ...string) {
	key := keys[0]
	if len(keys) == 1 {
		input[key] = value
	} else {
		if _, ok := input[key]; !ok {
			input[key] = Values{}
		}
		setNested(input[key].(Values), value, keys[1:]...)
	}
}

func (values Values) Val() any {
	return values["Value"]
}

func (values Values) V() any {
	return values["Value"]
}
