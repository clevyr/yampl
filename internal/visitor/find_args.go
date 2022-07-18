package visitor

import (
	"fmt"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"text/template/parse"
)

var fieldRe = regexp.MustCompile(`\.([A-Za-z_.]+)`)

func NewFindArgs(conf config.Config) FindArgs {
	return FindArgs{
		conf:     conf,
		matchMap: make(map[string]MatchSlice),
	}
}

type Match struct {
	Value    any
	Template string
	Position *token.Position
}

func (m Match) String() string {
	val := fmt.Sprintf("%v", m.Value)
	maxLen := 33
	if len(val) > maxLen {
		val = val[:maxLen-3] + "..."
	}
	var result string
	if m.Position != nil {
		result += "line " + strconv.Itoa(m.Position.Line) + ": "
	}
	result += fmt.Sprintf("%s %#v", val, m.Template)
	result = strings.ReplaceAll(result, "\n", " ")
	return result
}

type MatchSlice []Match

func (v MatchSlice) String() string {
	var s []string
	for _, match := range v {
		s = append(s, match.String())
	}
	return strings.Join(s, ";  ")
}

type FindArgs struct {
	conf     config.Config
	matchMap map[string]MatchSlice
	err      error
}

func (visitor *FindArgs) Visit(n ast.Node) ast.Visitor {
	if comment := node.GetCommentTmpl(visitor.conf.Prefix, n); comment != "" {
		tmpl, err := template.New("").
			Funcs(template2.FuncMap()).
			Delims(visitor.conf.LeftDelim, visitor.conf.RightDelim).
			Option("missingkey=zero").
			Parse(comment)
		if err != nil {
			visitor.err = err
			return nil
		}

		for _, field := range listTemplFields(tmpl) {
			if tokens := fieldRe.FindStringSubmatch(field); tokens != nil {
				for _, tok := range tokens[1:] {
					match := Match{
						Template: comment,
						Position: n.GetToken().Position,
					}
					switch n := n.(type) {
					case *ast.LiteralNode:
						match.Value = n.Value.String()
					default:
						if scalar, ok := n.(ast.ScalarNode); ok {
							match.Value = scalar.GetValue()
						}
					}
					visitor.matchMap[tok] = append(visitor.matchMap[tok], match)
				}
			}
		}
	}
	return visitor
}

func listTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}

func (visitor FindArgs) Values() []string {
	result := make([]string, 0, len(visitor.matchMap))
outer:
	for k, v := range visitor.matchMap {
		for kconf := range visitor.conf.Values {
			if k == kconf {
				continue outer
			}
		}
		for _, reserved := range config.ReservedKeys {
			if k == reserved {
				continue outer
			}
		}
		result = append(result, fmt.Sprintf("%s=\t%v", k, v))
	}
	return result
}

func (visitor FindArgs) Error() error {
	return visitor.err
}
