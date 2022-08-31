package visitor

import (
	"fmt"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"text/template/parse"
)

var fieldRe = regexp.MustCompile(`\.([A-Za-z\d_.]+)`)

func NewFindArgs(conf config.Config) FindArgs {
	return FindArgs{
		conf:     conf,
		matchMap: make(map[string]MatchSlice),
	}
}

type Match struct {
	Value    any
	Template string
	Line     int
	Column   int
}

func (m Match) String() string {
	val := fmt.Sprintf("%v", m.Value)
	maxLen := 33
	if len(val) > maxLen {
		val = val[:maxLen-3] + "..."
	}
	var result string
	if m.Line != 0 {
		result += "line " + strconv.Itoa(m.Line) + ": "
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
}

func (visitor *FindArgs) Visit(n *yaml.Node) error {
	if tmplSrc := node.GetCommentTmpl(visitor.conf.Prefix, n); tmplSrc != "" {
		tmpl, err := template.New("").
			Funcs(template2.FuncMap()).
			Delims(visitor.conf.LeftDelim, visitor.conf.RightDelim).
			Option("missingkey=zero").
			Parse(tmplSrc)
		if err != nil {
			return err
		}

		for _, field := range listTemplFields(tmpl) {
			if tokens := fieldRe.FindStringSubmatch(field); tokens != nil {
				for _, tok := range tokens[1:] {
					match := Match{
						Template: tmplSrc,
						Line:     n.Line,
						Column:   n.Column,
					}
					match.Value = n.Value
					visitor.matchMap[tok] = append(visitor.matchMap[tok], match)
				}
			}
		}
	}
	return nil
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
