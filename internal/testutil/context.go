package testutil

import (
	"context"
	"os"
	"stoke/internal/schema/policy"

	"github.com/rs/zerolog"
)

type ContextOption func(context.Context) context.Context

func NewMockContext(opts ...ContextOption) context.Context {
	ctx := policy.ConfigurePolicies([]string{}, []string{}, []string{}, "u", false, true, context.Background())
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
