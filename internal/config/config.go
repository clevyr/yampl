package config

type Config struct {
	Values     Values
	Inplace    bool
	Recursive  bool
	Prefix     string
	LeftDelim  string
	RightDelim string
	Indent     int
	Fail       bool
	Strip      bool
}

func New() *Config {
	return &Config{
		Values:     make(Values),
		Prefix:     "#yampl",
		LeftDelim:  "{{",
		RightDelim: "}}",
		Indent:     2,
	}
}
