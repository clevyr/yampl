package visitor

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
	template2 "github.com/clevyr/yampl/internal/template"
	"gopkg.in/yaml.v3"
)

func NewFindArgs(conf *config.Config) *FindArgs {
	return &FindArgs{
		conf:    conf,
		matches: make(map[string]MatchSlice),
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
		result += fmt.Sprintf("line %d: ", m.Line)
	}
	if val != "" {
		result += val + " "
	}
	result += fmt.Sprintf("%q", m.Template)
	result = strings.ReplaceAll(result, "\n", " ")
	return result
}

type MatchSlice []Match

func (v MatchSlice) String() string {
	s := make([]string, 0, len(v))
	for _, match := range v {
		s = append(s, match.String())
	}
	return strings.Join(s, "; ")
}

type FindArgs struct {
	conf    *config.Config
	path    string
	matches map[string]MatchSlice
}

func (f *FindArgs) Run(n *yaml.Node) error {
	switch {
	case len(n.Content) == 0:
		// Node has no children. Search current node.
		_ = f.FindArgs(n, n.Value)
	case n.Kind == yaml.MappingNode:
		for i := 0; i < len(n.Content); i += 2 {
			// Attempt to fetch template from comments on the key.
			key, val := n.Content[i], n.Content[i+1]

			tmplSrc, _ := comment.Parse(f.conf.Prefix, key)
			if tmplSrc == "" {
				// Key did not have comment, traversing children.
				if err := f.Run(val); err != nil {
					return err
				}
			} else {
				// Template is on key's comment instead of value.
				// This typically happens if the value is left empty with an implied null.
				_ = f.FindArgs(key, val.Value)
			}
		}
	default:
		for _, node := range n.Content {
			if err := f.Run(node); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *FindArgs) FindArgs(n *yaml.Node, value string) error {
	if tmplSrc, _ := comment.Parse(f.conf.Prefix, n); tmplSrc != "" {
		tmpl, err := template.New("").
			Funcs(template2.FuncMap()).
			Delims(f.conf.LeftDelim, f.conf.RightDelim).
			Option("missingkey=zero").
			Parse(tmplSrc)
		if err != nil {
			return NewNodeError(err, f.path, n)
		}

		for _, field := range listTmplFields(tmpl) {
			match := Match{
				Template: tmplSrc,
				Line:     n.Line,
				Column:   n.Column,
				Value:    value,
			}
			f.matches[field] = append(f.matches[field], match)
		}
	}
	return nil
}

func listTmplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	switch node := node.(type) {
	case *parse.ListNode:
		for _, node := range node.Nodes {
			res = listNodeFields(node, res)
		}
	case *parse.ActionNode:
		res = listNodeFields(node.Pipe, res)
	case *parse.PipeNode:
		for _, node := range node.Cmds {
			res = listNodeFields(node, res)
		}
	case *parse.CommandNode:
		for _, node := range node.Args {
			res = listNodeFields(node, res)
		}
	case *parse.FieldNode:
		res = append(res, node.Ident[0])
	}
	return res
}

func (f FindArgs) Values() []string {
	result := make([]string, 0, len(f.matches))
outer:
	for key, val := range f.matches {
		for kconf := range f.conf.Values {
			if key == kconf {
				continue outer
			}
		}
		for _, reserved := range []string{"Value", "Val", "V"} {
			if key == reserved {
				continue outer
			}
		}
		result = append(result, fmt.Sprintf("%s=\t%v", key, val))
	}
	return result
}
