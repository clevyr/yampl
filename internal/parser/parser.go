package parser

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"io"
)

func ParseBytes(b []byte) (*ast.File, error) {
	return parser.ParseBytes(append(b, '\n'), parser.ParseComments)
}

func ParseReader(r io.Reader) (*ast.File, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return ParseBytes(b)
}
