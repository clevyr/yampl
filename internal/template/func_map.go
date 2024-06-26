package template

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func FuncMap() template.FuncMap {
	fmap := sprig.TxtFuncMap()
	for k, v := range fmap {
		if strings.HasPrefix(k, "must") {
			k := strings.TrimPrefix(k, "must")
			k = strings.ToLower(k[0:1]) + k[1:]
			fmap[k] = v
		}
	}
	fmap["repo"] = DockerRepo
	fmap["tag"] = DockerTag
	fmap["toYaml"] = toYaml
	return fmap
}
