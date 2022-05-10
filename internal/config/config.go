package config

import (
	"io"
)

type Config struct {
	Paths      []string
	Values     Values
	Inline     bool
	Prefix     string
	LeftDelim  string
	RightDelim string
	Writer     io.Writer
}
