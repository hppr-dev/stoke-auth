package testutil

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type ContextOption func(context.Context) context.Context

func NewMockContext(opts ...ContextOption) context.Context {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return ctx
}

func StdLogger() ContextOption {
	return func(ctx context.Context) context.Context {
		return zerolog.New(os.Stdout).WithContext(ctx)
	}
}
