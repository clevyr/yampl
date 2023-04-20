package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues_Fill(t *testing.T) {
	type args struct {
		rawValues map[string]string
	}
	tests := []struct {
		name   string
		values Values
		args   args
		want   Values
	}{
		{"simple", make(Values), args{map[string]string{"a": "a"}}, Values{"a": "a"}},
		{"nested", make(Values), args{map[string]string{"a.b": "a"}}, Values{"a": Values{"b": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.values.Fill(tt.args.rawValues)
			assert.Equal(t, tt.want, tt.values)
		})
	}
}

func TestValues_SetNested(t *testing.T) {
	type args struct {
		v any
		k []string
	}
	tests := []struct {
		name   string
		values Values
		args   args
		want   Values
	}{
		{"simple", make(Values), args{"a", []string{"a"}}, Values{"a": "a"}},
		{"nested", make(Values), args{"a", []string{"a", "b"}}, Values{"a": Values{"b": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.values.SetNested(tt.args.v, tt.args.k...)
			assert.Equal(t, tt.want, tt.values)
		})
	}
}

func TestValues_V(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		want   any
	}{
		{"simple", Values{"Value": "a"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.values.V()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValues_Val(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		want   any
	}{
		{"simple", Values{"Value": "a"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.values.Val()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_setNested(t *testing.T) {
	type args struct {
		input Values
		value any
		keys  []string
	}
	tests := []struct {
		name string
		args args
		want Values
	}{
		{"simple", args{make(Values), "a", []string{"a"}}, Values{"a": "a"}},
		{"nested", args{make(Values), "a", []string{"a", "b"}}, Values{"a": Values{"b": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setNested(tt.args.input, tt.args.value, tt.args.keys...)
			assert.Equal(t, tt.want, tt.args.input)
		})
	}
}
