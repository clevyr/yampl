package node

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

func NewPrintableError(err error, node ast.Node) PrintableError {
	return PrintableError{
		err:  err,
		node: node,
	}
}

type PrintableError struct {
	err  error
	node ast.Node
}

func (p PrintableError) Error() string {
	return p.err.Error()
}

func (p PrintableError) AnnotateSource(src string, colored bool) string {
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
