package node

import (
	"gopkg.in/yaml.v3"
)

type Visitor func(node *yaml.Node) error

func Visit(visit Visitor, node *yaml.Node) error {
	if len(node.Content) == 0 {
		if err := visit(node); err != nil {
			return err
		}
	} else {
		for _, node := range node.Content {
			if err := Visit(visit, node); err != nil {
				return err
			}
		}
	}
	return nil
}
