package cfg

import (
	"context"
)

type ctxKey byte
const (
	configCtxKey ctxKey = iota
	serveMuxCtxKey
)

func Ctx(ctx context.Context) *Config {
	return ctx.Value(configCtxKey).(*Config)
}
