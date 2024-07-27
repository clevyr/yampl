package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVars_Fill(t *testing.T) {
	type args struct {
		rawValues map[string]string
	}
	tests := []struct {
		name string
		vars Vars
		args args
		want Vars
	}{
		{"simple", make(Vars), args{map[string]string{"a": "a"}}, Vars{"a": "a"}},
		{"nested", make(Vars), args{map[string]string{"a.b": "a"}}, Vars{"a": Vars{"b": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.vars.Fill(tt.args.rawValues)
			assert.Equal(t, tt.want, tt.vars)
		})
	}
}

func TestVars_SetNested(t *testing.T) {
	type args struct {
		v any
		k []string
	}
	tests := []struct {
		name string
		vars Vars
		args args
		want Vars
	}{
		{"simple", make(Vars), args{"a", []string{"a"}}, Vars{"a": "a"}},
		{"nested", make(Vars), args{"a", []string{"a", "b"}}, Vars{"a": Vars{"b": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.vars.SetNested(tt.args.v, tt.args.k...)
			assert.Equal(t, tt.want, tt.vars)
		})
	}
}

func TestVars_V(t *testing.T) {
	tests := []struct {
		name string
		vars Vars
		want any
	}{
		{"simple", Vars{"Value": "a"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.vars.V()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVars_Val(t *testing.T) {
	tests := []struct {
		name string
		vars Vars
		want any
	}{
		{"simple", Vars{"Value": "a"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.vars.Val()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_setNested(t *testing.T) {
	type args struct {
		input Vars
		value any
		keys  []string
	}
	tests := []struct {
		name string
		args args
		want Vars
	}{
		{"simple", args{make(Vars), "a", []string{"a"}}, Vars{"a": "a"}},
		{"nested", args{make(Vars), "a", []string{"a", "b"}}, Vars{"a": Vars{"b": "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setNested(tt.args.input, tt.args.value, tt.args.keys...)
			assert.Equal(t, tt.want, tt.args.input)
		})
	}
}
