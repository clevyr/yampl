package config

type Config struct {
	Values     Values
	Inline     bool
	Prefix     string
	LeftDelim  string
	RightDelim string
}

func New() Config {
	return Config{
		Values:     make(Values),
		Prefix:     "#yampl",
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
}
