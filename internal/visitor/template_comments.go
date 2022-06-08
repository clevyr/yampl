package visitor

import (
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	"github.com/goccy/go-yaml/ast"
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
	if v.err == nil {
		if err := node.Template(v.conf, n); err != nil {
			v.err = node.NewPrintableError(err, n)
			return nil
		}
	}
	return v
}

func (v TemplateComments) Error() error {
	return v.err
}
