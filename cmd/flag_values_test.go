package cmd

import (
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

func TestValMap_Slice(t *testing.T) {
	conf.Values["b"] = "b"

	tests := []struct {
		name string
		v    ValMap
		want []string
	}{
		{"simple", ValMap{"a": struct{}{}}, []string{"a="}},
		{"nested", ValMap{"a.b": struct{}{}}, []string{"a.b="}},
		{"duplicate", ValMap{"b": struct{}{}}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Slice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listNodeFields(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"simple", args{"{{ .a }}"}, []string{"{{.a}}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, _ := template.New("").Parse(tt.args.source)

			if got := listNodeFields(tmpl.Tree.Root, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listNodeFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listTemplFields(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"simple", args{"{{ .a }}"}, []string{"{{.a}}"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, _ := template.New("").Parse(tt.args.source)

			if got := listTemplFields(tmpl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listTemplFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValMap_Visitor(t *testing.T) {
	defaultConf := config.Config{
		LeftDelim:  "{{",
		RightDelim: "}}",
		Prefix:     "#yampl",
		Values: config.Values{
			"b": "b",
		},
	}

	tests := []struct {
		name    string
		v       ValMap
		source  string
		want    ValMap
		wantErr bool
	}{
		{"simple", make(ValMap), "a #yampl {{ .a }}", ValMap{"a": struct{}{}}, false},
		{"invalid template", make(ValMap), "a #yampl {{", ValMap{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := yaml.NewDecoder(strings.NewReader(tt.source))
			var n yaml.Node
			_ = decoder.Decode(&n)

			visitor := tt.v.Visitor()
			if err := visitor(defaultConf, n.Content[0]); (err != nil) != tt.wantErr {
				t.Errorf("Visitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.v, tt.want) {
				t.Errorf("Visitor() = %v, want %v", tt.v, tt.want)
			}
		})
	}
}
