package config

import "context"

type contextKey uint8

const configCtx contextKey = iota

func WithContext(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, configCtx, conf)
}

func FromContext(ctx context.Context) (*Config, bool) {
	conf, ok := ctx.Value(configCtx).(*Config)
	return conf, ok
}
