package cmd

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
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
		{"default", args{"color"}, &log.TextFormatter{}},
		{"plain", args{"plain"}, &log.TextFormatter{DisableColors: true}},
		{"json", args{"json"}, &log.JSONFormatter{}},
		{"unknown", args{""}, &log.TextFormatter{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initLogFormat(tt.args.format); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initLogFormat() = %v, want %v", got, tt.want)
			}
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
			if got := initLogLevel(tt.args.level); got != tt.want {
				t.Errorf("initLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
