package visitor

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func NewNodeError(err error, name string, node *yaml.Node) NodeError {
	return NodeError{
		err:  err,
		name: name,
		node: node,
	}
}

type NodeError struct {
	err  error
	name string
	node *yaml.Node
}

func (e NodeError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.name, e.node.Line, e.node.Column, e.err)
}

func (e NodeError) Unwrap() error {
	return e.err
}
