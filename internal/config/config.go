package config

import "github.com/rs/zerolog"

type Config struct {
	Values     Values
	Inplace    bool
	Recursive  bool
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
		Values:     make(Values),
		Prefix:     "#yampl",
		LeftDelim:  "{{",
		RightDelim: "}}",
		Indent:     2,

		IgnoreUnsetErrors: true,

		LogLevel:  zerolog.InfoLevel.String(),
		LogFormat: Auto,
	}
}
