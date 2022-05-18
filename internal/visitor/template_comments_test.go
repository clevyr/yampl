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
			file, _ := parser.ParseBytes([]byte("a " + tt.args.comment))
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
