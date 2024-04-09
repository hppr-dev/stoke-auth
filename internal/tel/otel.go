package tel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("stoke", trace.WithInstrumentationVersion("0.1"))
}

func GetMeter() metric.Meter {
	return otel.GetMeterProvider().Meter("stoke", metric.WithInstrumentationVersion("0.1"))
}

