package visitor

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type NodeErr struct {
	Err  error
	Node *yaml.Node
}

func (e NodeErr) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Node.Line, e.Node.Column, e.Err)
}
