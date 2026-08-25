package comment

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

// Parse returns the template, tag, and matched comment token from an ast.Node comment.
func Parse(prefix string, n ast.Node) (string, Tag, *token.Token) {
	return ParseGroup(prefix, n.GetComment())
}

// ParseGroup returns the template, tag, and matched comment token from a comment group.
func ParseGroup(prefix string, comments *ast.CommentGroupNode) (string, Tag, *token.Token) {
	if comments == nil {
		return "", NoTag, nil
	}

	for _, comment := range comments.Comments {
		s := comment.String()
		if strings.HasPrefix(s, prefix) {
			tk := comment.GetToken()
			// Comment has #yampl prefix
			s = strings.TrimPrefix(s, prefix)

			if strings.HasPrefix(s, " ") {
				// Tag not provided
				return s[1:], NoTag, tk
			}

			if strings.HasPrefix(s, tagSep) {
				// Match comment tag
				s = strings.TrimPrefix(s, tagSep)

				for _, tag := range Tags() {
					prefix := string(tag) + " "
					if strings.HasPrefix(s, prefix) {
						return s[len(prefix):], tag, tk
					}
				}
			}
		}
	}
	return "", NoTag, nil
}
