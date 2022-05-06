package config

type Values map[string]string

func (values Values) Val() string {
	return values["Value"]
}

func (values Values) V() string {
	return values["Value"]
}
