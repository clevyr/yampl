package template

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/clevyr/yampl/internal/util"
)

func FuncMap(opts ...Option) template.FuncMap {
	funcMap := sprig.TxtFuncMap()

	// Prefix functions with "may" that have a "must" counterpart
	for key, fn := range funcMap {
		if !strings.HasPrefix(key, "must") {
			mustKey := "must" + util.UpperFirst(key)
			if _, ok := funcMap[mustKey]; ok {
				key = "may" + util.UpperFirst(key)
				funcMap[key] = fn
			}
		}
	}

	// Remove prefix from "must" functions
	for key, fn := range funcMap {
		if strings.HasPrefix(key, "must") {
			k := strings.TrimPrefix(key, "must")
			k = util.LowerFirst(k)
			funcMap[k] = fn
		}
	}

	funcMap["repo"] = DockerRepo
	funcMap["tag"] = DockerTag
	funcMap["toYaml"] = toYaml

	for _, opt := range opts {
		opt(funcMap)
	}

	return funcMap
}
