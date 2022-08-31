package visitor

import (
	"bytes"
	"fmt"
	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/node"
	template2 "github.com/clevyr/yampl/internal/template"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"text/template"
)

func NewTemplateComments(conf config.Config) TemplateComments {
	return TemplateComments{
		conf: conf,
	}
}

type TemplateComments struct {
	conf config.Config
}

func (t TemplateComments) Visit(n *yaml.Node) error {
	if tmplSrc, tmplTag := node.GetCommentTmpl(t.conf.Prefix, n); tmplSrc != "" {
		t.conf.Log = t.conf.Log.WithFields(log.Fields{
			"tmpl":    tmplSrc,
			"filePos": fmt.Sprintf("%d:%d", n.Line, n.Column),
			"from":    n.Value,
		})

		tmpl, err := template.New("").
			Funcs(template2.FuncMap()).
			Delims(t.conf.LeftDelim, t.conf.RightDelim).
			Option("missingkey=error").
			Parse(tmplSrc)
		if err != nil {
			if !t.conf.Fail {
				t.conf.Log.WithError(err).Warn("skipping value due to template error")
				return nil
			}
			return NodeErr{Err: err, Node: n}
		}

		if t.conf.Values != nil {
			t.conf.Values["Value"] = n.Value
		}

		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, t.conf.Values); err != nil {
			if !t.conf.Fail {
				t.conf.Log.WithError(err).Warn("skipping value due to template error")
				return nil
			}
			return NodeErr{Err: err, Node: n}
		}

		if buf.String() != n.Value {
			t.conf.Log.WithField("to", buf.String()).Debug("updating value")
			n.Style = 0

			switch tmplTag {
			case node.SeqTag, node.MapTag:
				var tmpNode yaml.Node

				if err := yaml.Unmarshal(buf.Bytes(), &tmpNode); err != nil {
					if !t.conf.Fail {
						t.conf.Log.WithError(err).Warn("skipping value due to unmarshal error")
						return nil
					}
					return NodeErr{Err: err, Node: n}
				}

				content := tmpNode.Content[0]
				n.Content = content.Content
				n.Kind = content.Kind
				n.Value = content.Value
			default:
				n.SetString(buf.String())
			}

			n.Tag = tmplTag.ToYaml()
		}
	}
	return nil
}
