package comment

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse returns the template and tag from a yaml.Node LineComment.
func Parse(prefix string, n *yaml.Node) (string, Tag) {
	comment := n.LineComment
	if strings.HasPrefix(comment, prefix) {
		// Comment has #yampl prefix
		comment = strings.TrimPrefix(comment, prefix)

		if strings.HasPrefix(comment, " ") {
			// Tag not provided
			return comment[1:], DynamicTag
		}

		if strings.HasPrefix(comment, tagSep) {
			// Match comment tag
			comment = strings.TrimPrefix(comment, tagSep)

			for _, tag := range Tags() {
				prefix := string(tag) + " "
				if strings.HasPrefix(comment, prefix) {
					return comment[len(prefix):], tag
				}
			}
		}
	}
	return "", DynamicTag
}

// Move moves a comment between yaml.Node entries after a style change.
// When a yaml.MappingNode or yaml.SequenceNode has an inline comment,
// the decoder will set LineComment differently according to the node's style.
//
// When value(s) are on a single line (flow style), LineComment will be set on the value.
// When value(s) are on multiple lines (block style), LineComment will be set on the key.
//
// If templating changes the node's style, the comment needs to move or else
// encoding errors will occur.
func Move(key, val *yaml.Node) {
	if val.Kind != yaml.SequenceNode && val.Kind != yaml.MappingNode {
		return
	}

	if len(val.Content) > 0 && val.LineComment != "" && key.LineComment == "" {
		// Flow to block style: move comment from value to key.
		key.LineComment = val.LineComment
		val.LineComment = ""
	} else if len(val.Content) == 0 && key.LineComment != "" && val.LineComment == "" {
		// Block to flow style: move comment from key to value.
		val.LineComment = key.LineComment
		key.LineComment = ""
	}
}
