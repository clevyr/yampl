package config

import (
	"strings"
)

// Deprecated: Use `current` template function instead.
const CurrentValueKey = "Value"

type Vars map[string]any

func (vars Vars) Fill(src map[string]string) {
	for key, val := range src {
		split := strings.Split(key, ".")
		if len(split) == 1 {
			vars[key] = val
		} else {
			vars.SetNested(val, split...)
		}
	}
}

func (vars Vars) SetNested(v any, k ...string) {
	setNested(vars, v, k...)
}

func setNested(vars Vars, val any, keys ...string) {
	key := keys[0]
	if len(keys) == 1 {
		vars[key] = val
	} else {
		if _, ok := vars[key]; !ok {
			vars[key] = Vars{}
		}
		setNested(vars[key].(Vars), val, keys[1:]...) //nolint:errcheck
	}
}

// Deprecated: Use `current` template function instead.
func (vars Vars) Value() any {
	return vars["Value"]
}

// Deprecated: Use `current` template function instead.
func (vars Vars) Val() any {
	if v, ok := vars["Val"]; ok {
		return v
	}
	return vars.Value()
}

// Deprecated: Use `current` template function instead.
func (vars Vars) V() any {
	if v, ok := vars["V"]; ok {
		return v
	}
	return vars.Value()
}
