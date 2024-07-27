package template

import "text/template"

type Option func(funcMap template.FuncMap)

func WithCurrent(v string) Option {
	return func(funcMap template.FuncMap) {
		funcMap["current"] = func() string { return v }
	}
}
