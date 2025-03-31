package visitor

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
	yamplTemplate "github.com/clevyr/yampl/internal/template"
	"gopkg.in/yaml.v3"
)

func NewTemplateComments(conf *config.Config, path string) TemplateComments {
	logger := slog.Default()
	if path != "" {
		logger = logger.With("file", path)
	}

	return TemplateComments{
		conf: conf,
		log:  logger,
		path: path,
	}
}

type TemplateComments struct {
	conf *config.Config
	log  *slog.Logger
	path string
}

func (t TemplateComments) Run(n *yaml.Node) error { //nolint:gocognit
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
	log := t.log.With(
		"tmpl", tmplSrc,
		"filePos", fmt.Sprintf("%d:%d", n.Line, n.Column),
		"from", n.Value,
	)

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
		//nolint:staticcheck
		if _, ok := data[config.CurrentValueKey]; !ok {
			data[config.CurrentValueKey] = n.Value
		}
		t.checkDeprecated(tmplSrc)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return NewNodeError(err, name, n)
	}

	str := buf.String()
	if str != n.Value {
		log.Debug("Updating value", "to", str)
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
			n.SetString(str)
			switch {
			case n.Style != yaml.LiteralStyle && !utf8.ValidString(str),
				strings.HasPrefix(str, "\t"):
				n.Style = yaml.DoubleQuotedStyle
			}
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
				slog.Warn(key + " is deprecated, use `current` instead")
			}
		}
	}
}

func (t TemplateComments) handleTemplateError(err error) error {
	level := slog.LevelWarn
	switch {
	case err != nil && strings.Contains(err.Error(), "map has no entry for key"):
		if t.conf.IgnoreUnsetErrors {
			level = slog.LevelDebug
		} else {
			return err
		}
	case t.conf.IgnoreTemplateErrors:
	default:
		return err
	}
	t.log.Log(context.Background(), level, "Skipping value due to template error", "error", err)
	return nil
}
