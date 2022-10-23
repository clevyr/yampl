package template

import "testing"

func Test_toYaml(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"map", args{map[string]any{"a": "a"}}, "a: a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toYaml(tt.args.v); got != tt.want {
				t.Errorf("toYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}
