package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_toYaml(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"map", args{map[string]any{"a": "b"}}, "a: b", require.NoError},
		{"slice", args{[]string{"a", "b"}}, "- a\n- b", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toYaml(tt.args.v)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
