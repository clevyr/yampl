package template

import (
	"bytes"
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"testing"
)

var defaultConf = config.Config{
	LeftDelim:  "{{",
	RightDelim: "}}",
	Prefix:     "#yampl",
	Values: map[string]string{
		"b": "b",
	},
}

func TestRecurseNode(t *testing.T) {
	type args struct {
		conf  config.Config
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"no comment", args{defaultConf, "a: a"}, "a: a", false},
		{"simple comment", args{defaultConf, "a: a #yampl b"}, "a: b #yampl b", false},
		{"dynamic comment", args{defaultConf, "a: a #yampl {{ .b }}"}, "a: b #yampl {{ .b }}", false},
		{"invalid template", args{defaultConf, "a: a #yampl {{"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node yaml.Node
			_ = yaml.Unmarshal([]byte(tt.args.input), &node)

			if err := RecurseNode(tt.args.conf, &node); err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("RecurseNode() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			got, _ := yaml.Marshal(&node)
			got = bytes.TrimRight(got, "\n")
			if string(got) != tt.want {
				t.Errorf("RecurseNode() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestTemplateLineComment(t *testing.T) {
	type args struct {
		conf    config.Config
		comment string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"no comment", args{defaultConf, ""}, "a", false},
		{"simple comment", args{defaultConf, "#yampl b"}, "b #yampl b", false},
		{"dynamic comment", args{defaultConf, "#yampl {{ .b }}"}, "b #yampl {{ .b }}", false},
		{"prefix", args{config.Config{Prefix: "#tmpl"}, "#tmpl b"}, "b #tmpl b", false},
		{"delimiters", args{config.Config{LeftDelim: "<{", RightDelim: "}>", Prefix: "#yampl"}, `#yampl <{ "b" }>`}, `b #yampl <{ "b" }>`, false},
		{"invalid template", args{defaultConf, "#yampl {{"}, "", true},
		{"invalid variable", args{defaultConf, "#yampl {{ .z }}"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := yaml.Node{
				Kind:        yaml.ScalarNode,
				Tag:         "!!str",
				Value:       "a",
				LineComment: tt.args.comment,
			}

			if err := LineComment(tt.args.conf, &node); err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("TemplateLineComment() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			got, _ := yaml.Marshal(&node)
			got = bytes.TrimRight(got, "\n")
			if string(got) != tt.want {
				t.Errorf("TemplateLineComment() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
