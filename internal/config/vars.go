package config

import (
	"log/slog"
	"strings"
)

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
const currentValueKey = "Value"

// injectedValue marks the implicit current value stored by InjectCurrent so
// the deprecated accessors can tell it apart from a user-defined var that
// happens to be named "Value".
//
// Deprecated: Use `current` template function instead.
type injectedValue string

// InjectCurrent stores the node's current value under currentValueKey unless
// the user defined their own var with that name.
//
// Deprecated: Use `current` template function instead.
func (vars Vars) InjectCurrent(val string) {
	if _, ok := vars[currentValueKey]; !ok {
		vars[currentValueKey] = injectedValue(val)
	}
}

// Deprecated: Use `current` template function instead.
func (vars Vars) current() any {
	if v, ok := vars[currentValueKey].(injectedValue); ok {
		return string(v)
	}
	return vars[currentValueKey]
}

func warnDeprecated(name string) {
	slog.Warn("`" + name + "`" + " is deprecated, use `current` instead")
}

// Deprecated: Use `current` template function instead.
func (vars Vars) Value() any {
	if v, ok := vars[currentValueKey]; ok {
		if inj, ok := v.(injectedValue); ok {
			warnDeprecated(".Value")
			return string(inj)
		}
		// User-defined var that happens to be named "Value"
		return v
	}
	warnDeprecated(".Value")
	return nil
}

// Deprecated: Use `current` template function instead.
func (vars Vars) Val() any {
	if v, ok := vars["Val"]; ok {
		return v
	}
	warnDeprecated(".Val")
	return vars.current()
}

// Deprecated: Use `current` template function instead.
func (vars Vars) V() any {
	if v, ok := vars["V"]; ok {
		return v
	}
	warnDeprecated(".V")
	return vars.current()
}
