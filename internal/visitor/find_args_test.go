package visitor

import (
	"strings"
	"testing"

	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/parser"
	"github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindArgs_Visit(t *testing.T) {
	confWithVar := config.New()
	confWithVar.Vars = config.Vars{"b": "b"}

	tests := []struct {
		name   string
		conf   *config.Config
		source string
		want   []string
	}{
		{"no template", config.New(), "a: a", nil},
		{"simple", config.New(), "a: a #yampl {{ .b }}", []string{"b"}},
		{"key comment", config.New(), "a: #yampl {{ .b }}", []string{"b"}},
		{"seq item", config.New(), "a:\n  - a #yampl {{ .b }}", []string{"b"}},
		{"anchor", config.New(), "a: &x a #yampl {{ .b }}", []string{"b"}},
		{"empty flow map", config.New(), "a: {} #yampl {{ .b }}", []string{"b"}},
		{"empty flow seq", config.New(), "a: [] #yampl {{ .b }}", []string{"b"}},
		{"multiple", config.New(), "a: a #yampl {{ .b }}{{ .c }}", []string{"b", "c"}},
		{"already set", confWithVar, "a: a #yampl {{ .b }}", nil},
		{"no fields", config.New(), "a: a #yampl {{ current }}", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.source))
			require.NoError(t, err)

			v := NewFindArgs(tt.conf)
			for _, doc := range file.Docs {
				ast.Walk(v, doc.Body)
			}

			got := make([]string, 0, len(tt.want))
			for _, val := range v.Values() {
				key, _, _ := strings.Cut(val, "=")
				got = append(got, key)
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
