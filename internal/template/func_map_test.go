package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncMap(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name     string
		args     args
		wantFunc string
	}{
		{"repo", args{}, "repo"},
		{"tag", args{}, "tag"},
		{"current", args{[]Option{WithCurrent("test")}}, "current"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FuncMap(tt.args.opts...)
			_, ok := got[tt.wantFunc]
			assert.True(t, ok)
		})
	}
}
