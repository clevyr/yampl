package template

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"strings"
	"text/template"
)

func RecurseNode(conf config.Config, node *yaml.Node) error {
	if len(node.Content) == 0 {
		if node.LineComment != "" && strings.HasPrefix(node.LineComment, conf.Prefix) {
			if err := TemplateLineComment(conf, node); err != nil {
				return err
			}
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

func TemplateLineComment(conf config.Config, node *yaml.Node) error {
	tmpl, err := template.New("").
		Funcs(sprig.TxtFuncMap()).
		Delims(conf.LeftDelim, conf.RightDelim).
		Option("missingkey=error").
		Parse(strings.TrimSpace(node.LineComment[len(conf.Prefix):]))
	if err != nil {
		return err
	}

	var buf strings.Builder
	if err = tmpl.Execute(&buf, conf.Values); err != nil {
		return err
	}

	node.Value = buf.String()
	return nil
}
