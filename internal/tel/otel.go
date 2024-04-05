package tel

import (
	"context"
	"stoke/internal/ctx"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OTEL struct {
	Context    *ctx.Context
	Tracer     *trace.TracerProvider
	Meter      *metric.MeterProvider
	propagator propagation.TextMapPropagator
	serverInfo *resource.Resource
}

func (o *OTEL) Init() error {

	if err := o.buildResource(); err != nil {
		return err
	}
	if err := o.buildPropagator(); err != nil {
		return err
	}
	if err := o.buildTracerProvider(); err != nil {
		return err
	}
	if err := o.buildMeterProvider(); err != nil {
		return err
	}

	otel.SetTextMapPropagator(o.propagator)
	otel.SetTracerProvider(o.Tracer)
	otel.SetMeterProvider(o.Meter)

	return nil
}

func (o *OTEL) buildPropagator() error {
	o.propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	return nil
}

func (o *OTEL) buildTracerProvider() error {
	traceExporter, err := otlptracegrpc.New(o.Context.AppContext)
	if err != nil {
		return err
	}
	o.Tracer = trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
		trace.WithBatchTimeout(time.Second)),
		trace.WithResource(o.serverInfo),
	)
	o.Context.OnShutdown(wrapFunc(o.Tracer.Shutdown))
	return nil
}

func (o *OTEL) buildMeterProvider() error {
	metricExporter, err := prometheus.New()
	if err != nil {
		return err
	}
	o.Meter = metric.NewMeterProvider(
		metric.WithReader(metricExporter),
		metric.WithResource(o.serverInfo),
	)
	o.Context.OnShutdown(wrapFunc(o.Meter.Shutdown))
	return nil
}

func (o *OTEL) buildResource() error {
	var err error
	o.serverInfo, err = resource.New(o.Context.AppContext,
		resource.WithContainerID(),
		resource.WithAttributes(attribute.String("service.name", "stoke")),
	)
	return err
}

func wrapFunc(f func(context.Context) error) ctx.ShutdownFunc {
	return func (c ctx.Context) error {
		return f(c.AppContext)
	}
}
