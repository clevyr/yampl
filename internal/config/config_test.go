package config

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
				Indent:     2,
				Log:        log.NewEntry(log.StandardLogger()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			assert.Equal(t, tt.want, got)
		})
	}
}
