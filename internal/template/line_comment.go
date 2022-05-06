package template

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"strings"
	"text/template"
)

var funcMap = sprig.TxtFuncMap()

func init() {
	funcMap["repo"] = DockerRepo
	funcMap["tag"] = DockerTag
}

type Visitor func(conf config.Config, node *yaml.Node) error

func VisitNodes(conf config.Config, visit Visitor, node *yaml.Node) error {
	if len(node.Content) == 0 {
		if err := visit(conf, node); err != nil {
			return err
		}
	} else {
		for _, node := range node.Content {
			if err := VisitNodes(conf, visit, node); err != nil {
				return err
			}
		}
	}
	return nil
}

func LineComment(conf config.Config, node *yaml.Node) error {
	if node.LineComment != "" && strings.HasPrefix(node.LineComment, conf.Prefix) {
		tmpl, err := template.New("").
			Funcs(funcMap).
			Delims(conf.LeftDelim, conf.RightDelim).
			Option("missingkey=error").
			Parse(strings.TrimSpace(node.LineComment[len(conf.Prefix):]))
		if err != nil {
			return err
		}

		if conf.Values != nil {
			conf.Values["Value"] = node.Value
		}

		var buf strings.Builder
		if err = tmpl.Execute(&buf, conf.Values); err != nil {
			return err
		}

		node.Value = buf.String()
	}
	return nil
}
