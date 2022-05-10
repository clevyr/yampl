package config

import (
	"reflect"
	"testing"
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
			if !reflect.DeepEqual(tt.values, tt.want) {
				t.Errorf("Fill() = %v, want %v", tt.values, tt.want)
			}
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
			if !reflect.DeepEqual(tt.values, tt.want) {
				t.Errorf("SetNested() = %v, want %v", tt.values, tt.want)
			}
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
			if got := tt.values.V(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("V() = %v, want %v", got, tt.want)
			}
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
			if got := tt.values.Val(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Val() = %v, want %v", got, tt.want)
			}
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
			if !reflect.DeepEqual(tt.args.input, tt.want) {
				t.Errorf("setNested() = %v, want %v", tt.args.input, tt.want)
			}
		})
	}
}
