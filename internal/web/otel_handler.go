package web

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)


func WithSpan(label string, h http.Handler) spanHandler {
	return spanHandler{
		inner : otelhttp.NewHandler(h, label, otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
		label : label,
	}
}

type spanHandler struct {
	inner http.Handler
	label string
}

func (s spanHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	_ , span := trace.SpanFromContext(ctx).TracerProvider().Tracer("web-tracer").Start(ctx, s.label)
	defer span.End()

	otelhttp.WithRouteTag("inner"+s.label, s.inner).ServeHTTP(res, req)
}
