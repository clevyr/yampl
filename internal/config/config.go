package config

import (
	"log/slog"
	"strings"

	"github.com/clevyr/yampl/internal/config/flag"
)

type Config struct {
	valuesStringToString *flag.StringToString
	Vars                 Vars

	Inplace         bool
	Prefix          string
	LeftDelim       string
	RightDelim      string
	Indent          int
	Strip           bool
	NoSourceComment bool

	IgnoreUnsetErrors    bool
	IgnoreTemplateErrors bool

	LogLevel  string
	LogFormat string

	Completion string
}

func New() *Config {
	return &Config{
		valuesStringToString: &flag.StringToString{},
		Vars:                 make(Vars),

		Prefix:     "#yampl",
		LeftDelim:  "{{",
		RightDelim: "}}",
		Indent:     2,

		IgnoreUnsetErrors: true,

		LogLevel:  strings.ToLower(slog.LevelInfo.String()),
		LogFormat: FormatAuto.String(),
	}
}
