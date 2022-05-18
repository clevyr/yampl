package template

import (
	"github.com/Masterminds/sprig/v3"
	"text/template"
)

func FuncMap() template.FuncMap {
	fmap := sprig.TxtFuncMap()
	fmap["repo"] = DockerRepo
	fmap["tag"] = DockerTag
	return fmap
}
