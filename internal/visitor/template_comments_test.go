package visitor

import (
	"errors"
	"github.com/clevyr/go-yampl/internal/config"
	"github.com/clevyr/go-yampl/internal/parser"
	"reflect"
	"strings"
	"testing"
)

func TestNewTemplateComments(t *testing.T) {
	type args struct {
		conf config.Config
	}
	tests := []struct {
		name string
		args args
		want TemplateComments
	}{
		{"default", args{conf: config.Config{}}, TemplateComments{conf: config.Config{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTemplateComments(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTemplateComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateComments_Error(t *testing.T) {
	type fields struct {
		conf config.Config
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "no error", fields: fields{}},
		{"error", fields{err: errors.New("error")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := TemplateComments{
				conf: tt.fields.conf,
				err:  tt.fields.err,
			}
			if err := v.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplateComments_Visit(t *testing.T) {
	defaultConf := config.New()
	defaultConf.Values["b"] = "b"

	failConf := config.New()
	failConf.Fail = true

	prefixConf := config.New()
	prefixConf.Prefix = "#tmpl"

	delimConf := config.New()
	delimConf.LeftDelim = "<{"
	delimConf.RightDelim = "}>"
	delimConf.Prefix = "#yampl"

	stripConf := config.New()
	stripConf.Strip = true

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
		{"no comment", args{defaultConf, "test: a"}, "test: a", false},
		{"simple comment", args{defaultConf, "test: a #yampl b"}, "test: b #yampl b", false},
		{"dynamic comment", args{defaultConf, "test: a #yampl {{ .b }}"}, "test: b #yampl {{ .b }}", false},
		{"prefix", args{prefixConf, "test: a #tmpl b"}, "test: b #tmpl b", false},
		{"to string", args{prefixConf, "test: 1 #tmpl a"}, "test: a #tmpl a", false},
		{"to int", args{prefixConf, "test: a #tmpl 1"}, "test: 1 #tmpl 1", false},
		{"to float", args{prefixConf, "test: a #tmpl 0.1"}, "test: 0.1 #tmpl 0.1", false},
		{"same value", args{prefixConf, "test: a #tmpl a"}, "test: a #tmpl a", false},
		{"delimiters", args{delimConf, `test: a #yampl <{ "b" }>`}, `test: b #yampl <{ "b" }>`, false},
		{"sequence", args{delimConf, `- a #yampl b`}, `- b #yampl b`, false},
		{"invalid template", args{defaultConf, "test: a #yampl {{"}, "", true},
		{"mapping invalid variable ignore", args{defaultConf, "test: a #yampl {{ .z }}"}, "test: a #yampl {{ .z }}", false},
		{"sequence invalid variable ignore", args{defaultConf, "- a #yampl {{ .z }}"}, "- a #yampl {{ .z }}", false},
		{"mapping invalid variable error", args{failConf, "test: a #yampl {{ .z }}"}, "", true},
		{"sequence invalid variable error", args{failConf, "- a #yampl {{ .z }}"}, "", true},
		{"strip", args{stripConf, "test: a #yampl b"}, "test: b", false},
		{"no value", args{defaultConf, "test: #yampl a"}, "test: a #yampl a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := parser.ParseBytes([]byte(tt.args.comment))
			v := TemplateComments{conf: tt.args.conf}

			if v.Visit(file.Docs[0].Body); v.err != nil {
				if (v.err != nil) != tt.wantErr {
					t.Errorf("TemplateLineComment() error = %v, wantErr %v", v.err, tt.wantErr)
				}
				return
			}

			got := file.Docs[0].Body.String()
			got = strings.TrimRight(got, "\n")
			if got != tt.want {
				t.Errorf("TemplateLineComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
