package node

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type TmplTag string

func (t TmplTag) ToYaml() string {
	if t == DynamicTag {
		return ""
	}
	return "!!" + string(t)
}

const tagSep = ":"

var (
	DynamicTag TmplTag = ""
	BoolTag    TmplTag = "bool"
	StrTag     TmplTag = "str"
	IntTag     TmplTag = "int"
	FloatTag   TmplTag = "float"
	SeqTag     TmplTag = "seq"
	MapTag     TmplTag = "map"
)

var tags = []TmplTag{
	BoolTag,
	StrTag,
	IntTag,
	FloatTag,
	SeqTag,
	MapTag,
}

// GetCommentTmpl returns the template and tag from a yaml.Node LineComment
func GetCommentTmpl(prefix string, n *yaml.Node) (string, TmplTag) {
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

			for _, tag := range tags {
				prefix := string(tag) + " "
				if strings.HasPrefix(comment, prefix) {
					return comment[len(prefix):], tag
				}
			}
		}
	}
	return "", DynamicTag
}
