package config

import log "github.com/sirupsen/logrus"

type Config struct {
	Values     Values
	Inplace    bool
	Prefix     string
	LeftDelim  string
	RightDelim string
	Indent     int
	Fail       bool
	Log        *log.Entry
}

func New() Config {
	return Config{
		Values:     make(Values),
		Prefix:     "#yampl",
		LeftDelim:  "{{",
		RightDelim: "}}",
		Indent:     2,
		Log:        log.NewEntry(log.StandardLogger()),
	}
}
