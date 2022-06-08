package visitor

import (
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"github.com/goccy/go-yaml/ast"
	"regexp"
	"text/template"
	"text/template/parse"
)

var fieldRe = regexp.MustCompile(`\.([A-Za-z_.]+)`)

func NewFindArgs(conf config.Config) FindArgs {
	return FindArgs{
		conf:   conf,
		valMap: make(map[string]struct{}),
	}
}

type FindArgs struct {
	conf   config.Config
	valMap map[string]struct{}
	err    error
}

func (v *FindArgs) Visit(n ast.Node) ast.Visitor {
	if comment := node.GetCommentTmpl(v.conf.Prefix, n); comment != "" {
		tmpl, err := template.New("").
			Funcs(template2.FuncMap()).
			Delims(v.conf.LeftDelim, v.conf.RightDelim).
			Option("missingkey=zero").
			Parse(comment)
		if err != nil {
			v.err = err
			return nil
		}

		for _, field := range listTemplFields(tmpl) {
			matches := fieldRe.FindStringSubmatch(field)
			if matches != nil {
				for _, match := range matches[1:] {
					v.valMap[match] = struct{}{}
				}
			}
		}
	}
	return v
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

func (v FindArgs) Values() []string {
	result := make([]string, 0, len(v.valMap))
outer:
	for k := range v.valMap {
		for kconf := range v.conf.Values {
			if k == kconf {
				continue outer
			}
		}
		for _, reserved := range config.ReservedKeys {
			if k == reserved {
				continue outer
			}
		}
		result = append(result, k+"=")
	}
	return result
}

func (v FindArgs) Error() error {
	return v.err
}
