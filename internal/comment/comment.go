package comment

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
)

// Parse returns the template and tag from an ast.Node comment.
func Parse(prefix string, n ast.Node) (string, Tag) {
	return ParseGroup(prefix, n.GetComment())
}

// ParseGroup returns the template and tag from a comment group.
func ParseGroup(prefix string, comments *ast.CommentGroupNode) (string, Tag) {
	if comments == nil {
		return "", NoTag
	}

	for _, comment := range comments.Comments {
		s := comment.String()
		if strings.HasPrefix(s, prefix) {
			// Comment has #yampl prefix
			s = strings.TrimPrefix(s, prefix)

			if strings.HasPrefix(s, " ") {
				// Tag not provided
				return s[1:], NoTag
			}

			if strings.HasPrefix(s, tagSep) {
				// Match comment tag
				s = strings.TrimPrefix(s, tagSep)

				for _, tag := range Tags() {
					prefix := string(tag) + " "
					if strings.HasPrefix(s, prefix) {
						return s[len(prefix):], tag
					}
				}
			}
		}
	}
	return "", NoTag
}
