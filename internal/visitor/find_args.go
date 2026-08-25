package visitor

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/clevyr/yampl/internal/comment"
	"github.com/clevyr/yampl/internal/config"
	yamplTemplate "github.com/clevyr/yampl/internal/template"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
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
	matches map[string]MatchSlice
}

//nolint:ireturn
func (f *FindArgs) Visit(n ast.Node) ast.Visitor {
	tmplSrc, _, _ := comment.Parse(f.conf.Prefix, n)
	if tmplSrc == "" {
		// Comments on empty flow collections attach to the closing
		// bracket's next token instead of the node itself
		if val, ok := n.(*ast.MappingNode); ok && val.End != nil && val.End.NextType() == token.CommentType {
			tmplSrc, _, _ = comment.ParseGroup(f.conf.Prefix, ast.CommentGroup([]*token.Token{val.End.Next}))
		}
	}

	if tmplSrc != "" {
		f.FindArgs(n, tmplSrc)
	}
	return f
}

func (f *FindArgs) FindArgs(n ast.Node, tmplSrc string) {
	value := nodeValue(n)

	tmpl, err := template.New("").
		Funcs(yamplTemplate.FuncMap(
			yamplTemplate.WithCurrent(value),
		)).
		Delims(f.conf.LeftDelim, f.conf.RightDelim).
		Option("missingkey=zero").
		Parse(tmplSrc)
	if err != nil {
		return
	}

	for _, field := range listTmplFields(tmpl) {
		match := Match{
			Template: tmplSrc,
			Value:    value,
		}
		if pos := n.GetToken().Position; pos != nil {
			match.Line = pos.Line
			match.Column = pos.Column
		}
		f.matches[field] = append(f.matches[field], match)
	}
}

func nodeValue(n ast.Node) string {
	if scalar, ok := n.(ast.ScalarNode); ok {
		return fmt.Sprintf("%v", scalar.GetValue())
	}
	return ""
}

func listTmplFields(t *template.Template) []string {
	return listNodeFields(t.Root, nil)
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

func (f *FindArgs) Values() []string {
	result := make([]string, 0, len(f.matches))
outer:
	for key, val := range f.matches {
		for kconf := range f.conf.Vars {
			if key == kconf {
				continue outer
			}
		}
		result = append(result, fmt.Sprintf("%s=\t%v", key, val))
	}
	return result
}
