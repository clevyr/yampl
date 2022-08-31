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

func GetCommentTmpl(prefix string, n *yaml.Node) (string, TmplTag) {
	for _, comment := range []string{n.LineComment, n.HeadComment, n.FootComment} {
		fullPrefix := prefix + " "
		if strings.HasPrefix(comment, fullPrefix) {
			return comment[len(fullPrefix):], DynamicTag
		}
		for _, tag := range tags {
			fullPrefix := prefix + ":" + string(tag) + " "
			if strings.HasPrefix(comment, fullPrefix) {
				return comment[len(fullPrefix):], tag
			}
		}
	}
	return "", DynamicTag
}
