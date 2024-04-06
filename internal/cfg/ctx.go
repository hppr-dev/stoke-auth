package cfg

import (
	"context"

	"github.com/rs/zerolog"
)

func Ctx(ctx context.Context) *Config {
	return ctx.Value("config").(*Config)
}

func augmentContext(ctx context.Context, componentName string) context.Context {
	rootLogger := zerolog.Ctx(ctx)
	logger := rootLogger.With().
		Str("component", componentName).
		Logger()
	return logger.WithContext(ctx)
}
