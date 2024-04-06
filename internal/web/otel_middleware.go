package web

import (
	"net/http"
	"stoke/internal/tel"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)


func TraceHTTP(h http.Handler) spanHandler {
	return spanHandler{
		inner : otelhttp.NewHandler(h, "HTTP", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
		tracer: tel.GetTracer(),
	}
}

type spanHandler struct {
	inner http.Handler
	tracer trace.Tracer
}

func (s spanHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	label := req.URL.String()

	ctx, span := s.tracer.Start(ctx, label,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.ClientAddress(req.RemoteAddr),
				semconv.UserAgentOriginal(req.UserAgent()),
				semconv.TLSEstablished(req.TLS == nil),
				semconv.NetworkProtocolName(req.Proto),
				attribute.String("http.request.method", req.Method),
		),
	)
	defer span.End()


	otelhttp.WithRouteTag(label, s.inner).ServeHTTP(res, req.WithContext(ctx))
}
