package visitor

import (
	"bytes"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
	log "github.com/sirupsen/logrus"
	"text/template"
)

func NewTemplateComments(conf config.Config) TemplateComments {
	return TemplateComments{
		conf: conf,
	}
}

type TemplateComments struct {
	conf config.Config
	err  error
}

func (v *TemplateComments) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.MappingValueNode:
		if comment := node.GetCommentTmpl(v.conf.Prefix, n.Value); comment != "" {
			newNode, err := templateComment(v.conf, comment, n.Value)
			if err != nil {
				v.err = err
				return nil
			}

			if newNode != nil {
				if err := n.Replace(newNode); err != nil {
					v.err = err
					return nil
				}
			}
		}
	case *ast.SequenceNode:
		for i, value := range n.Values {
			if comment := node.GetCommentTmpl(v.conf.Prefix, value); comment != "" {
				newNode, err := templateComment(v.conf, comment, value)
				if err != nil {
					v.err = err
					return nil
				}

				if newNode != nil {
					if err := n.Replace(i, newNode); err != nil {
						v.err = err
						return nil
					}
				}
			}
		}
	}
	return v
}

func (v TemplateComments) Error() error {
	return v.err
}

func templateComment(conf config.Config, comment string, n ast.Node) (ast.Node, error) {
	tmpl, err := template.New("").
		Funcs(template2.FuncMap()).
		Delims(conf.LeftDelim, conf.RightDelim).
		Option("missingkey=error").
		Parse(comment)
	if err != nil {
		return nil, err
	}

	if conf.Values != nil {
		conf.Values["Value"] = n.(ast.ScalarNode).GetValue()
	}

	logEntry := conf.Log.WithField("yamlpath", n.GetPath())

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, conf.Values); err != nil {
		if !conf.Strict {
			logEntry.WithError(err).Warn("skipping value due to template error")
			return nil, nil
		}
		return nil, err
	}

	oldVal := n.GetToken().Value
	if buf.String() != oldVal {
		logEntry.WithFields(log.Fields{
			"tmpl": comment,
			"from": oldVal,
			"to":   buf.String(),
		}).Debug("updating value")

		tok := token.New(buf.String(), n.GetToken().Origin, n.GetToken().Position)
		var newNode ast.Node
		switch tok.Type {
		case token.IntegerType, token.BinaryIntegerType, token.OctetIntegerType, token.HexIntegerType:
			newNode = ast.Integer(tok)
		case token.FloatType:
			newNode = ast.Float(tok)
		default:
			newNode = ast.String(tok)
		}

		if err := newNode.SetComment(n.GetComment()); err != nil {
			return newNode, err
		}

		return newNode, nil
	}
	return nil, nil
}
