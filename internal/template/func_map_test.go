package template

import (
	"testing"
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
			if _, ok := got[tt.wantFunc]; !ok {
				t.Errorf("FuncMap() func not set: %v", tt.wantFunc)
			}
		})
	}
}
