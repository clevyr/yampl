package cmd

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_initLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"default"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initLog()
		})
	}
}

func Test_initLogFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name string
		args args
		want log.Formatter
	}{
		{"default", args{"auto"}, &log.TextFormatter{}},
		{"color", args{"color"}, &log.TextFormatter{ForceColors: true}},
		{"plain", args{"plain"}, &log.TextFormatter{DisableColors: true}},
		{"json", args{"json"}, &log.JSONFormatter{}},
		{"unknown", args{""}, &log.TextFormatter{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := initLogFormat(tt.args.format)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_initLogLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want log.Level
	}{
		{"trace", args{"trace"}, log.TraceLevel},
		{"debug", args{"debug"}, log.DebugLevel},
		{"info", args{"info"}, log.InfoLevel},
		{"warning", args{"warning"}, log.WarnLevel},
		{"error", args{"error"}, log.ErrorLevel},
		{"fatal", args{"fatal"}, log.FatalLevel},
		{"panic", args{"panic"}, log.PanicLevel},
		{"unknown", args{""}, log.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := initLogLevel(tt.args.level)
			assert.Equal(t, tt.want, got)
		})
	}
}
