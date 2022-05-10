package node

import (
	"errors"
	"github.com/clevyr/go-yampl/internal/config"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestVisitNodes(t *testing.T) {
	defaultConf := config.Config{
		LeftDelim:  "{{",
		RightDelim: "}}",
		Prefix:     "#yampl",
		Values: config.Values{
			"b": "b",
		},
	}

	type args struct {
		conf  config.Config
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{defaultConf, "a: a"}, false},
		{"error", args{defaultConf, "a: a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var visitorCalled bool

			var node yaml.Node
			_ = yaml.Unmarshal([]byte(tt.args.input), &node)

			visitor := func(conf config.Config, node *yaml.Node) error {
				visitorCalled = true
				if tt.wantErr {
					return errors.New("test error")
				}
				return nil
			}

			if err := Visit(tt.args.conf, visitor, &node); err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("Visit() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if !visitorCalled {
				t.Errorf("Visit() visitorCalled = %v, want %v", visitorCalled, true)
			}
		})
	}
}
