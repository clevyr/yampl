package config

type Config struct {
	Values     Values
	Inplace    bool
	Prefix     string
	LeftDelim  string
	RightDelim string
	Indent     int
	Strict     bool
}

func New() Config {
	return Config{
		Values:     make(Values),
		Prefix:     "#yampl",
		LeftDelim:  "{{",
		RightDelim: "}}",
		Indent:     2,
	}
}
