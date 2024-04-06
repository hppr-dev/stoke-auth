package web

import (
	"context"
	"net/http"
	"stoke/internal/ent"
	"stoke/internal/ent/ogent"

	"github.com/rs/zerolog"
)

func NewEntityAPIHandler(prefix string, ctx context.Context) http.Handler {
	if len(prefix) > 0 && prefix[len(prefix) - 1] == '/' {
		prefix = prefix[0:len(prefix) - 1]
	}

	hdlr, err := ogent.NewServer(
		ogent.NewOgentHandler(ent.FromContext(ctx)),
		ogent.WithPathPrefix(prefix),
	)
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("An error occured creating an entity handler")
		return nil
	}

	// Unwrap ServeHTTP to be have tracing spans use context from our custom middleware
	return http.HandlerFunc(hdlr.ServeHTTP)
}
