package visitor

import (
	"bytes"
	"fmt"
	"maps"
	"regexp"
	"strings"
	"text/template"

	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
	yamplTemplate "github.com/clevyr/yampl/internal/template"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func NewTemplateComments(conf *config.Config, path string) TemplateComments {
	var l zerolog.Logger
	if path == "" {
		l = log.Logger
	} else {
		l = log.With().Str("file", path).Logger()
	}

	return TemplateComments{
		conf: conf,
		log:  l,
		path: path,
	}
}

type TemplateComments struct {
	conf *config.Config
	log  zerolog.Logger
	path string
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

			if err := t.Template(t.path, n, tmplSrc, tmplTag); err != nil {
				if err := t.handleTemplateError(err); err != nil {
					return err
				}
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

				if err := t.Template(t.path, val, tmplSrc, tmplTag); err != nil {
					if err := t.handleTemplateError(err); err != nil {
						return err
					}
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

func (t TemplateComments) Template(name string, n *yaml.Node, tmplSrc string, tmplTag comment.Tag) error {
	log := t.log.With().
		Str("tmpl", tmplSrc).
		Str("filePos", fmt.Sprintf("%d:%d", n.Line, n.Column)).
		Str("from", n.Value).
		Logger()

	tmpl, err := template.New("`"+tmplSrc+"`").
		Funcs(yamplTemplate.FuncMap(
			yamplTemplate.WithCurrent(n.Value),
		)).
		Delims(t.conf.LeftDelim, t.conf.RightDelim).
		Option("missingkey=error").
		Parse(tmplSrc)
	if err != nil {
		return NewNodeError(err, name, n)
	}

	data := maps.Clone(t.conf.Vars)
	if data != nil {
		if _, ok := data[config.CurrentValueKey]; !ok {
			data[config.CurrentValueKey] = n.Value
		}
		t.checkDeprecated(tmplSrc)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return NewNodeError(err, name, n)
	}

	if buf.String() != n.Value {
		log.Debug().Str("to", buf.String()).Msg("updating value")
		n.Style = 0

		switch tmplTag {
		case comment.SeqTag, comment.MapTag:
			var tmpNode yaml.Node

			if err := yaml.Unmarshal(buf.Bytes(), &tmpNode); err != nil {
				return NewNodeError(err, name, n)
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

func (t TemplateComments) checkDeprecated(tmplSrc string) {
	if t.conf.Vars != nil {
		re := regexp.MustCompile(`(\.V(al(ue)?)?)(?:[ |)]|` + regexp.QuoteMeta(t.conf.RightDelim) + `)`)
		for _, match := range re.FindAllStringSubmatch(tmplSrc, -1) {
			key := match[1]
			if _, ok := t.conf.Vars[key[1:]]; !ok {
				log.Warn().Msg(key + " is deprecated, use `current` instead")
			}
		}
	}
}

func (t TemplateComments) handleTemplateError(err error) error {
	level := zerolog.WarnLevel
	switch {
	case err != nil && strings.Contains(err.Error(), "map has no entry for key"):
		if t.conf.IgnoreUnsetErrors {
			level = zerolog.DebugLevel
		} else {
			return err
		}
	case t.conf.IgnoreTemplateErrors:
	default:
		return err
	}
	t.log.WithLevel(level).Err(err).Msg("skipping value due to template error")
	return nil
}
