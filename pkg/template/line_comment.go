package template

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"strings"
	"text/template"
)

type LineComment struct {
	Config config.Config
	root   *yaml.Node
}

func (l *LineComment) UnmarshalYAML(node *yaml.Node) error {
	l.root = node
	return recurseNode(l, node)
}

func recurseNode(l *LineComment, node *yaml.Node) error {
	if len(node.Content) == 0 {
		if node.LineComment != "" && strings.HasPrefix(node.LineComment, l.Config.Prefix) {
			if err := templateLineComment(l, node); err != nil {
				return err
			}
		}
	} else {
		for _, node := range node.Content {
			if err := recurseNode(l, node); err != nil {
				return err
			}
		}
	}
	return nil
}

func templateLineComment(l *LineComment, node *yaml.Node) error {
	tmpl, err := template.New("").
		Funcs(sprig.TxtFuncMap()).
		Delims(l.Config.LeftDelim, l.Config.RightDelim).
		Option("missingkey=error").
		Parse(strings.TrimSpace(node.LineComment[len(l.Config.Prefix):]))
	if err != nil {
		return err
	}
	var buf strings.Builder
	if err = tmpl.Execute(&buf, l.Config.Values); err != nil {
		return err
	}
	node.Value = buf.String()
	return nil
}

func (l LineComment) MarshalYAML() (any, error) {
	return l.root, nil
}
