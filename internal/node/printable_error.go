package node

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
)

func NewPrintableError(err error, node ast.Node) PrintableError {
	return PrintableError{
		err:  err,
		node: node,
	}
}

// NewPrintableTemplateError returns a PrintableError that annotates the
// template inside a comment instead of the node's value.
func NewPrintableTemplateError(err error, node ast.Node, tmpl string, tk *token.Token) PrintableError {
	return PrintableError{
		err:  err,
		node: node,
		tmpl: tmpl,
		tk:   tk,
	}
}

type PrintableError struct {
	err  error
	node ast.Node
	tmpl string
	tk   *token.Token
}

func (p PrintableError) Error() string {
	return p.err.Error()
}

func (p PrintableError) Unwrap() error {
	return p.err
}

func (p PrintableError) AnnotateSource(src string, colored bool) string {
	if p.tk != nil {
		return p.annotateTemplate(colored)
	}

	path, err := yaml.PathString(p.node.GetPath())
	if err != nil {
		return ""
	}

	source, err := path.AnnotateSource([]byte(src), colored)
	if err != nil {
		return ""
	}

	return string(source)
}

// annotateTemplate annotates the template inside a comment. The token is cloned
// so that shifting the annotation off of the comment prefix does not move the
// comment itself.
func (p PrintableError) annotateTemplate(colored bool) string {
	tk := p.tk.Clone()
	tk.Position.Column += len(tk.Value) - len(p.tmpl) + len("#")

	var pp printer.Printer
	return pp.PrintErrorToken(tk, colored)
}
