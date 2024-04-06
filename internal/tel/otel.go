package tel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("stoke", trace.WithInstrumentationVersion("0.1"))
}

