package cmd

import (
	"errors"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/node"
	template2 "github.com/clevyr/go-yampl/internal/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
	"text/template/parse"
)

func init() {
	Command.Flags().StringToStringVarP((*map[string]string)(&conf.Values), "value", "v", map[string]string{}, "Define a template variable. Can be used more than once.")
	err := Command.RegisterFlagCompletionFunc("value", valueCompletion)
	if err != nil {
		panic(err)
	}
}

func valueCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	valMap := make(ValMap)
	visitor := valMap.Visitor()

	for _, path := range args {
		func() {
			f, err := os.Open(path)
			if err != nil {
				return
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			decoder := yaml.NewDecoder(f)

			var n yaml.Node
			if err := decoder.Decode(&n); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
			}

			_ = node.Visit(conf, visitor, &n)
		}()
	}

	return valMap.Slice(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
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

type ValMap map[string]struct{}

func (v ValMap) Slice() []string {
	result := make([]string, 0, len(v))
outer:
	for k := range v {
		for kconf := range conf.Values {
			if k == kconf {
				continue outer
			}
		}
		result = append(result, k+"=")
	}
	return result
}

func (v ValMap) Visitor() func(conf config.Config, node *yaml.Node) error {
	var fieldRe = regexp.MustCompile(`\.([A-Za-z_.]+)`)

	return func(conf config.Config, node *yaml.Node) error {
		if node.LineComment != "" && strings.HasPrefix(node.LineComment, conf.Prefix) {
			tmpl, err := template.New("").
				Funcs(template2.FuncMap).
				Delims(conf.LeftDelim, conf.RightDelim).
				Option("missingkey=zero").
				Parse(strings.TrimSpace(node.LineComment[len(conf.Prefix):]))
			if err != nil {
				return err
			}

			for _, field := range listTemplFields(tmpl) {
				matches := fieldRe.FindStringSubmatch(field)
				if matches != nil {
					for _, match := range matches[1:] {
						v[match] = struct{}{}
					}
				}
			}
		}
		return nil
	}
}
