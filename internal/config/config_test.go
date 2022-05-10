package config

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			"defaults",
			Config{
				Values:     make(Values),
				Prefix:     "#yampl",
				LeftDelim:  "{{",
				RightDelim: "}}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
