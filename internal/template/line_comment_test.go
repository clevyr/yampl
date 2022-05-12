package template

import (
	"bytes"
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestTemplateLineComment(t *testing.T) {
	defaultConf := config.New()
	defaultConf.Values["b"] = "b"

	strictConf := config.New()
	strictConf.Strict = true

	prefixConf := config.New()
	prefixConf.Prefix = "#tmpl"

	delimConf := config.New()
	delimConf.LeftDelim = "<{"
	delimConf.RightDelim = "}>"
	delimConf.Prefix = "#yampl"

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
		{"prefix", args{prefixConf, "#tmpl b"}, "b #tmpl b", false},
		{"delimiters", args{delimConf, `#yampl <{ "b" }>`}, `b #yampl <{ "b" }>`, false},
		{"invalid template", args{defaultConf, "#yampl {{"}, "", true},
		{"invalid variable ignore", args{defaultConf, "#yampl {{ .z }}"}, "a #yampl {{ .z }}", false},
		{"invalid variable error", args{strictConf, "#yampl {{ .z }}"}, "", true},
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
				t.Errorf("TemplateLineComment() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}
