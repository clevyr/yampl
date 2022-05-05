package config

type Config struct {
	Paths      []string
	Values     map[string]string
	Inline     bool
	Prefix     string
	LeftDelim  string
	RightDelim string
}
