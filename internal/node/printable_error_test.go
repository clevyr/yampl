package node

import (
	"errors"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// findComment returns the comment token holding tmpl, mimicking comment.Parse.
func findComment(t *testing.T, src, tmpl string) *token.Token {
	t.Helper()

	for _, tk := range lexer.Tokenize(src) {
		if tk.Type == token.CommentType && strings.HasSuffix(tk.Value, tmpl) {
			return tk
		}
	}
	return nil
}

func TestPrintableError_AnnotateSource(t *testing.T) {
	tests := []struct {
		name string
		src  string
		path string
		tmpl string
		want string
	}{
		{
			"value comment",
			"a: \"\" #yampl {{ .test\nb: \"\" #yampl {{ .ok }}\n",
			"$.a",
			"{{ .test",
			">  1 | a: \"\" #yampl {{ .test\n                    ^\n   2 | b: \"\" #yampl {{ .ok }}\n   3 | ",
		},
		{
			"tagged comment",
			"a: \"\" #yampl:int notanint\n",
			"$.a",
			"notanint",
			">  1 | a: \"\" #yampl:int notanint\n                        ^\n",
		},
		{
			"key comment",
			"a: #yampl {{ .test\n  b: c\n",
			"$.a",
			"{{ .test",
			">  1 | a: #yampl {{ .test\n                 ^\n   2 |   b: c",
		},
		{
			"sequence value",
			"a:\n  - \"\" #yampl {{ .test\n",
			"$.a[0]",
			"{{ .test",
			"   1 | a:\n>  2 |   - \"\" #yampl {{ .test\n                     ^\n",
		},
		{
			"annotates the value without a comment token",
			"a: b\n",
			"$.a",
			"",
			">  1 | a: b\n          ^\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(tt.src), parser.ParseComments)
			require.NoError(t, err)

			path, err := yaml.PathString(tt.path)
			require.NoError(t, err)

			n, err := path.FilterFile(file)
			require.NoError(t, err)

			tk := findComment(t, tt.src, tt.tmpl)
			p := NewPrintableTemplateError(errors.New("test"), n, tt.tmpl, tk)
			assert.Equal(t, tt.want, p.AnnotateSource(tt.src, false))

			if tk != nil {
				assert.Equal(t, tt.src, file.String(), "annotating must not modify the source")
			}
		})
	}
}
