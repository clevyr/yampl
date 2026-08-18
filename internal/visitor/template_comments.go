package visitor

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"maps"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/node"
	yamplTemplate "github.com/clevyr/yampl/internal/template"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
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

//nolint:ireturn
func (t *TemplateComments) Visit(n ast.Node) ast.Visitor {
	if t.err == nil {
		if err := t.Run(n); err != nil {
			t.err = node.NewPrintableError(err, n)
			return nil
		}
	}
	return t
}

type TemplateComments struct {
	conf *config.Config
	log  *slog.Logger
	path string
	err  error
}

func (t *TemplateComments) Error() error {
	return t.err
}

func (t *TemplateComments) Run(n ast.Node) error { //nolint:gocognit,gocyclo,cyclop
	switch n := n.(type) {
	case *ast.MappingValueNode:
		// Comments attach to the value an anchor wraps, not the anchor itself
		anchor, _ := n.Value.(*ast.AnchorNode)
		value := n.Value
		if anchor != nil {
			value = anchor.Value
		}

		// Comment after value
		tmplSrc, tmplTag := comment.Parse(t.conf.Prefix, value)
		if t.conf.Strip {
			if err := value.SetComment(nil); err != nil {
				return err
			}
		}

		if tmplSrc == "" { //nolint:nestif
			// Edge case where comment is set on key
			if tmplSrc, tmplTag = comment.Parse(t.conf.Prefix, n.Key); tmplSrc != "" {
				if !t.conf.Strip {
					// Move comment from key to value
					if err := value.SetComment(n.Key.GetComment()); err != nil {
						return err
					}
				}
				if err := n.Key.SetComment(nil); err != nil {
					return err
				}
			}

			if tmplSrc == "" {
				// Check for empty map with comment as next
				val, ok := value.(*ast.MappingNode)
				if ok && val.End != nil && val.End.NextType() == token.CommentType {
					if !t.conf.Strip {
						if err := val.SetComment(ast.CommentGroup([]*token.Token{val.End.Next})); err != nil {
							return err
						}
					}
					val.End.Next = nil
					tmplSrc, tmplTag = comment.Parse(t.conf.Prefix, val)
				}
			}
		}

		if tmplSrc != "" {
			newNode, err := t.Template(n.Key, value, tmplSrc, tmplTag)
			if err != nil {
				return t.handleTemplateError(n, tmplSrc, err)
			}

			if newNode != nil {
				if anchor != nil {
					anchor.Value = newNode
				} else if err := n.Replace(newNode); err != nil {
					return err
				}
			}
		}
	case *ast.SequenceNode:
		for i, value := range n.Values {
			anchor, _ := value.(*ast.AnchorNode)
			if anchor != nil {
				value = anchor.Value
			}

			tmplSrc, tmplTag := comment.Parse(t.conf.Prefix, value)
			if tmplSrc == "" {
				continue
			}

			if t.conf.Strip {
				if err := value.SetComment(nil); err != nil {
					return err
				}
			}

			newNode, err := t.Template(nil, value, tmplSrc, tmplTag)
			if err != nil {
				return t.handleTemplateError(value, tmplSrc, err)
			}

			if newNode != nil {
				if anchor != nil {
					anchor.Value = newNode
				} else if err := n.Replace(i, newNode); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//nolint:gocognit,gocyclo,cyclop,nestif,ireturn,funlen
func (t *TemplateComments) Template(
	kn ast.Node,
	vn ast.Node,
	tmplSrc string,
	tmplTag comment.Tag,
) (ast.Node, error) {
	oldVal := vn.GetToken().Value
	log := t.nodeLogger(vn, tmplSrc)

	tmpl, err := template.New("").
		Funcs(yamplTemplate.FuncMap(
			yamplTemplate.WithCurrent(oldVal),
		)).
		Delims(t.conf.LeftDelim, t.conf.RightDelim).
		Option("missingkey=error").
		Parse(tmplSrc)
	if err != nil {
		return nil, err
	}

	data := maps.Clone(t.conf.Vars)
	if data != nil {
		data.InjectCurrent(oldVal) //nolint:staticcheck // Supports the deprecated .Value var
	}

	var buf strings.Builder
	buf.Grow(len(oldVal))
	if err = tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	if buf.String() != oldVal {
		log.Debug("Updating value", "to", buf.String())

		var v any
		switch tmplTag {
		case comment.NoTag, comment.StrTag:
			s := buf.String()
			if !utf8.ValidString(s) {
				// Matches gopkg.in/yaml.v3, which base64-encodes invalid UTF-8
				s = base64.StdEncoding.EncodeToString([]byte(s))
			}
			v = s
		case comment.IntTag:
			v, err = strconv.ParseInt(buf.String(), 10, 64)
			if err != nil {
				return nil, err
			}
		case comment.FloatTag:
			v, err = strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return nil, err
			}
		case comment.BoolTag:
			v, err = strconv.ParseBool(buf.String())
			if err != nil {
				return nil, err
			}
		case comment.SeqTag:
			var seq []any
			if err := yaml.NewDecoder(strings.NewReader(buf.String())).Decode(&seq); err != nil {
				return nil, err
			}
			v = seq
		case comment.MapTag:
			var m map[any]any
			if err := yaml.NewDecoder(strings.NewReader(buf.String())).Decode(&m); err != nil {
				return nil, err
			}
			v = m
		}

		n, err := yaml.NewEncoder(
			nil,
			yaml.Indent(t.conf.Indent),
			yaml.IndentSequence(true),
			yaml.UseLiteralStyleIfMultiline(true),
		).EncodeToNode(v)
		if err != nil {
			return nil, err
		}

		if str, ok := n.(*ast.StringNode); ok {
			lbc := token.DetectLineBreakCharacter(str.Value)
			switch {
			case needsDoubleQuote(str.Value):
				str.Token.Type = token.DoubleQuoteType
			case strings.Contains(str.Value, lbc):
				// Prepare multiline string nodes for comment
				b, err := str.MarshalYAML()
				if err != nil {
					return nil, err
				}

				// Only use a block scalar if it round-trips losslessly
				var check string
				if err := yaml.Unmarshal(b, &check); err == nil && check == str.Value {
					if err := yaml.Unmarshal(b, &n); err != nil {
						return nil, err
					}
				} else {
					str.Token.Type = token.DoubleQuoteType
				}
			default:
				// goccy's needs-quotes check misses some plain scalars that
				// other parsers treat as non-strings or reject entirely
				// (e.g. "30E3" parses as a float, "?" as a complex key).
				// str.Value holds the rendered form, so only plain renderings
				// need the check.
				if str.Value == v && isAmbiguousPlain(str.Value) {
					str.Token.Type = token.DoubleQuoteType
				}
			}
		}

		isFlowStyle := getIsFlowStyle(n)
		switch n.(type) {
		case *ast.MappingNode, *ast.MappingValueNode, *ast.SequenceNode:
			// Children of a block collection render at the replaced value's
			// column. Anchor it to the key so the indent is independent of
			// where the old value happened to sit on the line.
			if !isFlowStyle {
				if kn != nil {
					vn.GetToken().Position.Column = kn.GetToken().Position.Column + t.conf.Indent
				} else {
					vn.GetToken().Position.Column--
				}
			}
		}

		if seq, ok := n.(*ast.SequenceNode); ok {
			if isFlowStyle {
				b, err := seq.MarshalYAML()
				if err != nil {
					return nil, err
				}

				if err := yaml.Unmarshal(b, &n); err != nil {
					return nil, err
				}

				if c := vn.GetComment(); c != nil {
					if err := n.SetComment(c); err != nil {
						return nil, err
					}

					n.GetToken().Next = token.Comment(c.String(), c.String(), &token.Position{})
				}
			} else if kn != nil {
				if err := kn.SetComment(vn.GetComment()); err != nil {
					return nil, err
				}
			}
		} else {
			if err := n.SetComment(vn.GetComment()); err != nil {
				return nil, err
			}
		}

		return n, nil
	}
	return nil, nil //nolint:nilnil // nil node means the value was unchanged
}

func (t *TemplateComments) nodeLogger(n ast.Node, tmplSrc string) *slog.Logger {
	oldVal := n.GetToken().Value

	log := t.log.With("path", n.GetPath(), "tmpl", tmplSrc)
	if pos := n.GetToken().Position; pos != nil {
		log = log.With("file_pos", fmt.Sprintf("%d:%d", pos.Line, pos.Column))
	}
	log = log.With("from", oldVal)
	return log
}

func (t *TemplateComments) handleTemplateError(n ast.Node, tmplSrc string, err error) error {
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

	t.nodeLogger(n, tmplSrc).Log(context.Background(), level,
		"Skipping value due to template error",
		"error", err,
	)
	return nil
}
