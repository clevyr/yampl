package visitor

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type NodeError struct {
	Err  error
	Name string
	Node *yaml.Node
}

func (e NodeError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.Name, e.Node.Line, e.Node.Column, e.Err)
}

func (e NodeError) Unwrap() error {
	return e.Err
}
