package node

import (
	"github.com/goccy/go-yaml/ast"
	"strings"
)

func GetCommentTmpl(prefix string, node ast.Node) string {
	comments := node.GetComment()
	if comments != nil {
		for _, comment := range comments.Comments {
			s := comment.String()
			if strings.HasPrefix(s, prefix) {
				s = strings.TrimPrefix(s, prefix)
				s = strings.TrimSpace(s)
				return s
			}
		}
	}
	return ""
}
