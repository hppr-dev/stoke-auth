package web

import (
	"net/http"
	"stoke/internal/ctx"
	"stoke/internal/ent/ogent"
	"stoke/internal/tel"
)

type EntityAPI struct {
	*ogent.OgentHandler
	Context *ctx.Context
}

func NewEntityAPIHandler(prefix string, context *ctx.Context, o *tel.OTEL) http.Handler {
	if len(prefix) > 0 && prefix[len(prefix) - 1] == '/' {
		prefix = prefix[0:len(prefix) - 1]
	}
	hdlr, err := ogent.NewServer(
		EntityAPI{
			OgentHandler: ogent.NewOgentHandler(context.DB),
			Context: context,
		},
		ogent.WithPathPrefix(prefix),
		ogent.WithTracerProvider(o.Tracer),
		ogent.WithMeterProvider(o.Meter),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("An error occured creating an entity handler")
		return nil
	}
	return hdlr

}
