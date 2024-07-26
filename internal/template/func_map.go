package template

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/clevyr/yampl/internal/util"
)

func FuncMap() template.FuncMap {
	fmap := sprig.TxtFuncMap()
	// Prefix functions with "may" that have a "must" counterpart
	for key, fn := range fmap {
		if !strings.HasPrefix(key, "must") {
			mustKey := "must" + util.UpperFirst(key)
			if _, ok := fmap[mustKey]; ok {
				key = "may" + util.UpperFirst(key)
				fmap[key] = fn
			}
		}
	}
	// Remove prefix from "must" functions
	for key, fn := range fmap {
		if strings.HasPrefix(key, "must") {
			k := strings.TrimPrefix(key, "must")
			k = util.LowerFirst(k)
			fmap[k] = fn
		}
	}
	fmap["repo"] = DockerRepo
	fmap["tag"] = DockerTag
	fmap["toYaml"] = toYaml
	return fmap
}
