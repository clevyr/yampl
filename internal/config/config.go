package config

import (
	"github.com/clevyr/yampl/internal/config/flag"
	"github.com/rs/zerolog"
)

type Config struct {
	valuesStringToString *flag.StringToString
	Vars                 Vars

	Inplace    bool
	Prefix     string
	LeftDelim  string
	RightDelim string
	Indent     int
	Strip      bool

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

		LogLevel:  zerolog.InfoLevel.String(),
		LogFormat: Auto,
	}
}
