package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToString_Set(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    *StringToString
		wantErr require.ErrorAssertionFunc
	}{
		{
			"one value",
			args{"a=b"},
			&StringToString{value: map[string]string{"a": "b"}, changed: true},
			require.NoError,
		},
		{
			"two values",
			args{"a=b,c=d"},
			&StringToString{value: map[string]string{"a": "b", "c": "d"}, changed: true},
			require.NoError,
		},
		{
			"multiline value",
			args{"a=b\nc,d=e"},
			&StringToString{value: map[string]string{"a": "b\nc", "d": "e"}, changed: true},
			require.NoError,
		},
		{
			"multiline values",
			args{"a=b\nc=d"},
			&StringToString{value: map[string]string{"a": "b", "c": "d"}, changed: true},
			require.NoError,
		},
		{
			"multiple newlines",
			args{"a=b\n\nc=d"},
			&StringToString{value: map[string]string{"a": "b", "c": "d"}, changed: true},
			require.NoError,
		},
		{
			"trim spaces",
			args{"a=b\n    c=d"},
			&StringToString{value: map[string]string{"a": "b", "c": "d"}, changed: true},
			require.NoError,
		},
		{
			"newline around values",
			args{"\na=b\nc=d\n"},
			&StringToString{value: map[string]string{"a": "b", "c": "d"}, changed: true},
			require.NoError,
		},
		{
			"json value",
			args{"a=[1]"},
			&StringToString{value: map[string]string{"a": "[1]"}, changed: true},
			require.NoError,
		},
		{
			"json values",
			args{"a=[1, 2, 3]"},
			&StringToString{value: map[string]string{"a": "[1, 2, 3]"}, changed: true},
			require.NoError,
		},
		{"error empty", args{""}, &StringToString{}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringToString{}
			tt.wantErr(t, s.Set(tt.args.val))
			assert.Equal(t, tt.want, s)
		})
	}

	t.Run("consecutive", func(t *testing.T) {
		s := &StringToString{}
		require.NoError(t, s.Set("a=b"))
		assert.True(t, s.changed)
		assert.Equal(t, map[string]string{"a": "b"}, s.value)
		require.NoError(t, s.Set("c=d"))
		assert.True(t, s.changed)
		assert.Equal(t, map[string]string{"a": "b", "c": "d"}, s.value)
	})
}

func TestStringToString_String(t *testing.T) {
	tests := []struct {
		name  string
		value *StringToString
		want  string
	}{
		{"empty", &StringToString{}, "[]"},
		{"simple value", &StringToString{value: map[string]string{"a": "b"}, changed: true}, "[a=b]"},
		{"value with comma", &StringToString{value: map[string]string{"a": "b,c"}, changed: true}, `["a=b,c"]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
