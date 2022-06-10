package node

import (
	"bytes"
	"github.com/clevyr/go-yampl/internal/config"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
	log "github.com/sirupsen/logrus"
	"text/template"
)

func Template(conf config.Config, n ast.Node) error {
	switch n := n.(type) {
	case *ast.MappingValueNode:
		// Comment after value
		comment := GetCommentTmpl(conf.Prefix, n.Value)

		// Edge case where comment is set on key
		if comment == "" {
			if comment = GetCommentTmpl(conf.Prefix, n.Key); comment != "" {
				// Move comment from key to value
				if err := n.Value.SetComment(n.Key.GetComment()); err != nil {
					return err
				}
				if err := n.Key.SetComment(nil); err != nil {
					return err
				}
			}
		}

		if comment != "" {
			newNode, err := templateComment(conf, comment, n.Value)
			if err != nil {
				return err
			}

			if newNode != nil {
				if err := n.Replace(newNode); err != nil {
					return err
				}
			}
		}
	case *ast.SequenceNode:
		for i, value := range n.Values {
			if comment := GetCommentTmpl(conf.Prefix, value); comment != "" {
				newNode, err := templateComment(conf, comment, value)
				if err != nil {
					return err
				}

				if newNode != nil {
					if err := n.Replace(i, newNode); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
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
		n, ok := n.(ast.ScalarNode)
		if ok {
			conf.Values["Value"] = n.GetValue()
		}
	}

	logEntry := conf.Log.WithField("yamlpath", n.GetPath())

	if conf.Strip {
		if err := n.SetComment(nil); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, conf.Values); err != nil {
		if !conf.Fail {
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
