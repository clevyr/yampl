package visitor

import (
	"regexp"
	"strings"

	"github.com/goccy/go-yaml/ast"
)

// needsDoubleQuote reports whether s can only be represented as a
// double-quoted scalar: it contains characters that are not YAML-printable
// (or that goccy renders raw), including tabs, which break plain scalars and
// block scalar indentation.
func needsDoubleQuote(s string) bool {
	for _, r := range s {
		switch {
		case r == '\n':
		case r == '\t', r < 0x20, r == 0x7F, 0x80 <= r && r <= 0x9F,
			r == 0x2028, r == 0x2029, r == 0xFEFF, r == 0xFFFE, r == 0xFFFF:
			return true
		}
	}
	return false
}

// numberLike matches strings that yaml.v3 (and other YAML 1.1-lineage
// parsers) resolve as numbers but goccy's needs-quotes check leaves plain:
// exponent forms without a dot ("30E3"), underscore digit groups ("1_000",
// "-_09"), ".inf"/".nan" variants, and base-prefixed ints ("0x83", "0X83")
// including yaml.v3's quirky embedded-sign forms ("0b-0").
//
// The rule set was derived from a differential sweep of goccy's rendering
// against yaml.v3's parser over an exhaustive corpus, and is continuously
// revalidated by Test_templateReader_differential and the fuzz test in
// internal/processor.
var numberLike = regexp.MustCompile(
	`^[-+]?(\.[0-9_]+|_*[0-9][0-9_]*(\.[0-9_]*)?)([eE][-+]?[0-9_]+)?$` +
		`|^[-+]?\.(inf|Inf|INF|nan|NaN|NAN)$` +
		`|^0[bBoOxX][-+0-9a-fA-F_]*$`,
)

// isAmbiguousPlain reports whether s, rendered as a plain scalar, would be
// misread by other YAML parsers as a non-string or rejected outright.
func isAmbiguousPlain(s string) bool {
	return numberLike.MatchString(s) || s == "?" || strings.HasPrefix(s, "? ")
}

func getIsFlowStyle(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SequenceNode:
		return n.IsFlowStyle || len(n.Values) == 0
	case *ast.MappingNode:
		return n.IsFlowStyle
	case *ast.MappingValueNode:
		return n.IsFlowStyle
	}
	return false
}
