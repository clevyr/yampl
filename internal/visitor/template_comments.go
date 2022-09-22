package visitor

import (
	"bytes"
	"fmt"
	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
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

func (t TemplateComments) Run(n *yaml.Node) error {
	if len(n.Content) == 0 {
		// Node has no children. Template current node.
		tmplSrc, tmplTag := comment.Parse(t.conf.Prefix, n)
		if tmplSrc != "" {
			if err := t.Template(n, tmplSrc, tmplTag); err != nil {
				if t.conf.Fail {
					return err
				} else {
					t.conf.Log.WithError(err).Warn("skipping value due to template error")
				}
			}
		}
		return nil
	}

	switch n.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(n.Content); i += 2 {
			// Attempt to fetch template from comments on the key.
			key, val := n.Content[i], n.Content[i+1]

			tmplSrc, tmplTag := comment.Parse(t.conf.Prefix, key)
			if tmplSrc != "" {
				if err := t.Template(val, tmplSrc, tmplTag); err != nil {
					if t.conf.Fail {
						return err
					} else {
						t.conf.Log.WithError(err).Warn("skipping value due to template error")
					}
				} else {
					// Current node was templated, do not need to traverse children
					comment.Move(key, val)
					continue
				}
			}

			// Key did not have comment, traversing children.
			if err := t.Run(val); err != nil {
				return err
			}

			comment.Move(key, val)
		}
	default:
		for _, n := range n.Content {
			if err := t.Run(n); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t TemplateComments) Template(n *yaml.Node, tmplSrc string, tmplTag comment.Tag) error {
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
		return NodeErr{Err: err, Node: n}
	}

	if t.conf.Values != nil {
		t.conf.Values["Value"] = n.Value
	}

	if t.conf.Strip {
		n.LineComment = ""
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, t.conf.Values); err != nil {
		return NodeErr{Err: err, Node: n}
	}

	if buf.String() != n.Value {
		t.conf.Log.WithField("to", buf.String()).Debug("updating value")
		n.Style = 0

		switch tmplTag {
		case comment.SeqTag, comment.MapTag:
			var tmpNode yaml.Node

			if err := yaml.Unmarshal(buf.Bytes(), &tmpNode); err != nil {
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
	return nil
}
