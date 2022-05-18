package visitor

import (
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"github.com/goccy/go-yaml/ast"
	log "github.com/sirupsen/logrus"
	"strings"
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
	if comment := node.GetCommentTmpl(v.conf.Prefix, n); comment != "" {
		tmpl, err := template.New("").
			Funcs(template2.FuncMap()).
			Delims(v.conf.LeftDelim, v.conf.RightDelim).
			Option("missingkey=error").
			Parse(comment)
		if err != nil {
			v.err = err
			return nil
		}

		if v.conf.Values != nil {
			v.conf.Values["Value"] = n.String()
		}

		l := v.conf.Log.WithField("yamlpath", n.(*ast.StringNode).BaseNode.Path)

		var buf strings.Builder
		if err = tmpl.Execute(&buf, v.conf.Values); err != nil {
			if !v.conf.Strict {
				l.WithError(err).Warn("skipping value due to template error")
				return nil
			}
			v.err = err
			return nil
		}

		if buf.String() != n.String() {
			l.WithFields(log.Fields{
				"tmpl": comment,
				"from": n.String(),
				"to":   buf.String(),
			}).Debug("updating value")
			n.(*ast.StringNode).Value = buf.String()
		}
	}
	return v
}

func (v TemplateComments) Error() error {
	return v.err
}
