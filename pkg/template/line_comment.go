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

func RecurseNode(conf config.Config, node *yaml.Node) error {
	if len(node.Content) == 0 {
		if err := LineComment(conf, node); err != nil {
			return err
		}
	} else {
		for _, node := range node.Content {
			if err := RecurseNode(conf, node); err != nil {
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
