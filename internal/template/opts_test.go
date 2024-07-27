package template

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithCurrent(t *testing.T) {
	funcMap := template.FuncMap{}
	WithCurrent("test")(funcMap)
	v, ok := funcMap["current"]
	require.True(t, ok, "current function should exist")
	fn, ok := v.(func() string)
	require.True(t, ok, "current function return string")
	got := fn()
	assert.Equal(t, "test", got)
}
