package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncMap(t *testing.T) {
	tests := []struct {
		name     string
		wantFunc string
	}{
		{"repo", "repo"},
		{"tag", "tag"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FuncMap()
			_, ok := got[tt.wantFunc]
			assert.True(t, ok)
		})
	}
}
