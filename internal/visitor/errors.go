package visitor

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type NodeError struct {
	Err  error
	Node *yaml.Node
}

func (e NodeError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Node.Line, e.Node.Column, e.Err)
}

func (e NodeError) Unwrap() error {
	return e.Err
}
