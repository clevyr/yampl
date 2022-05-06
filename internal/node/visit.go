package node

import (
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
)

type Visitor func(conf config.Config, node *yaml.Node) error

func Visit(conf config.Config, visit Visitor, node *yaml.Node) error {
	if len(node.Content) == 0 {
		if err := visit(conf, node); err != nil {
			return err
		}
	} else {
		for _, node := range node.Content {
			if err := Visit(conf, visit, node); err != nil {
				return err
			}
		}
	}
	return nil
}
