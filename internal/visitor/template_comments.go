package visitor

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
	template2 "github.com/clevyr/yampl/internal/template"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

func NewTemplateComments(conf *config.Config, log zerolog.Logger) TemplateComments {
	return TemplateComments{
		conf: conf,
		log:  log,
	}
}

type TemplateComments struct {
	conf *config.Config
	log  zerolog.Logger
}

func (t TemplateComments) Run(n *yaml.Node) error {
	switch {
	case len(n.Content) == 0:
		// Node has no children. Template current node.
		tmplSrc, tmplTag := comment.Parse(t.conf.Prefix, n)
		if tmplSrc != "" {
			if t.conf.Strip {
				n.LineComment = ""
			}

			if err := t.Template(n, tmplSrc, tmplTag); err != nil {
				if t.conf.Fail {
					return err
				}
				t.log.Warn().Err(err).Msg("skipping value due to template error")
			}
		}
	case n.Kind == yaml.MappingNode:
		for i := 0; i < len(n.Content); i += 2 {
			// Attempt to fetch template from comments on the key.
			key, val := n.Content[i], n.Content[i+1]

			tmplSrc, tmplTag := comment.Parse(t.conf.Prefix, key)
			if tmplSrc == "" {
				// Key did not have comment, traversing children.
				if err := t.Run(val); err != nil {
					return err
				}
			} else {
				// Template is on key's comment instead of value.
				// This typically happens if the value is left empty with an implied null.

				if t.conf.Strip {
					key.LineComment = ""
				}

				if err := t.Template(val, tmplSrc, tmplTag); err != nil {
					if t.conf.Fail {
						return err
					}
					t.log.Warn().Err(err).Msg("skipping value due to template error")
					continue
				}
			}
			comment.Move(key, val)
		}
	default:
		// Iterate over children
		for _, n := range n.Content {
			if err := t.Run(n); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t TemplateComments) Template(n *yaml.Node, tmplSrc string, tmplTag comment.Tag) error {
	log := t.log.With().
		Str("tmpl", tmplSrc).
		Str("filePos", fmt.Sprintf("%d:%d", n.Line, n.Column)).
		Str("from", n.Value).
		Logger()

	tmpl, err := template.New("").
		Funcs(template2.FuncMap()).
		Delims(t.conf.LeftDelim, t.conf.RightDelim).
		Option("missingkey=error").
		Parse(tmplSrc)
	if err != nil {
		return NodeError{Err: err, Node: n}
	}

	if t.conf.Values != nil {
		t.conf.Values["Value"] = n.Value
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, t.conf.Values); err != nil {
		return NodeError{Err: err, Node: n}
	}

	if buf.String() != n.Value {
		log.Debug().Str("to", buf.String()).Msg("updating value")
		n.Style = 0

		switch tmplTag {
		case comment.SeqTag, comment.MapTag:
			var tmpNode yaml.Node

			if err := yaml.Unmarshal(buf.Bytes(), &tmpNode); err != nil {
				return NodeError{Err: err, Node: n}
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
