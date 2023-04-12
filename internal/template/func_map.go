package template

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func FuncMap() template.FuncMap {
	fmap := sprig.TxtFuncMap()
	fmap["repo"] = DockerRepo
	fmap["tag"] = DockerTag
	fmap["toYaml"] = toYaml
	return fmap
}
