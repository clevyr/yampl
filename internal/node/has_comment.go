package node

import (
	"gopkg.in/yaml.v3"
	"strings"
)

func GetCommentTmpl(prefix string, n *yaml.Node) string {
	for _, comment := range []string{n.LineComment, n.HeadComment, n.FootComment} {
		if strings.HasPrefix(comment, prefix) {
			return strings.TrimSpace(comment[len(prefix):])
		}
	}
	return ""
}
