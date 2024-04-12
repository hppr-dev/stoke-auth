package cfg

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type ContextFunc func(context.Context) error

type Telemetry struct {
	// The type of metric provider. Can be http or grpc. Leave blank to disable.
	// Note: there will always be an internel prometheus metric exporter for monitoring from the UI
	MetricProvider  string `json:"metric_provider"`
	// The URL of where to export metrics.
	MetricExportURL string `json:"metric_exporter_url"`
	// Use an insecure connection to connect to MetricExportURL
	MetricInsecure  bool `json:"metric_insecure"`

	// The type of trace provider. Can be http or grpc. Leave blank to disable
	TraceProvider  string `json:"trace_provider"`
	// The URL of where to export traces
	TraceExportURL string `json:"trace_export_url"`
	// Use an insecure connection to connect to TraceExportURL
	TraceInsecure  bool `json:"trace_insecure"`

	// Disable exposing default prometheus metrics
	DisableDefaultPrometheus bool `json:"disable_default_prometheus"`
	// Require authentication for default prometheus
	RequirePrometheusAuthentication bool `json:"require_prometheus_authentication"`

	// Non-parsed fields
	shutdownFuncs []ContextFunc
}

func (t *Telemetry) Initialize(ctx context.Context) ([]ContextFunc, error) {

	info, err := t.buildResource(ctx)
	if err != nil {
		return nil, err
	}

	tracer, err := t.buildTracerProvider(info, ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(t.buildPropagator())
	if tracer != nil {
		otel.SetTracerProvider(tracer)
	}
	otel.SetMeterProvider(t.buildMeterProvider(info, ctx))

	return t.shutdownFuncs, nil
}

func (t *Telemetry) buildPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (t *Telemetry) buildTracerProvider(info *resource.Resource, ctx context.Context) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	var err error
	switch t.TraceProvider {
	case "grpc":
		var opts []otlptracegrpc.Option
		if t.TraceExportURL != "" {
		 opts = append(opts, otlptracegrpc.WithEndpointURL(t.TraceExportURL))
		}
		if t.TraceInsecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)

	case "http":
		var opts []otlptracehttp.Option
		if t.TraceExportURL != "" {
		 opts = append(opts, otlptracehttp.WithEndpointURL(t.TraceExportURL))
		}
		if t.TraceInsecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		exporter, err = otlptracehttp.New(ctx, opts...)

	default:
		return nil, nil
	}
	tracer := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(info),
	)
	t.shutdownFuncs = append(t.shutdownFuncs, tracer.Shutdown)
	return tracer, err
}

func (t *Telemetry) buildMeterProvider(info *resource.Resource, ctx context.Context) *metric.MeterProvider {
	var exporter metric.Exporter
	var err error

	readerOptions := []metric.Option{
		metric.WithResource(info),
	}

	switch t.MetricProvider {
	case "grpc":
		var opts []otlpmetricgrpc.Option
		if t.MetricExportURL != "" {
		 opts = append(opts, otlpmetricgrpc.WithEndpointURL(t.MetricExportURL))
		}
		if t.MetricInsecure {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		}
		exporter, err = otlpmetricgrpc.New(ctx, opts...)

	case "http":
		var opts []otlpmetrichttp.Option
		if t.MetricExportURL != "" {
		 opts = append(opts, otlpmetrichttp.WithEndpointURL(t.MetricExportURL))
		}
		if t.MetricInsecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}
		exporter, err = otlpmetrichttp.New(ctx, opts...)

	}

	if err != nil {
		zerolog.Ctx(ctx).Info().
			Err(err).
			Str("exporterType", t.MetricProvider).
			Msg("Could not initialize configured metric exporter")
	} else if exporter != nil {
		readerOptions = append(readerOptions, metric.WithReader(metric.NewPeriodicReader(exporter)))
	}

	promExporter, err := prometheus.New()
	if err != nil {
		zerolog.Ctx(ctx).Fatal().
			Err(err).
			Msg("Could not initialize prometheus exporter")
	}
	readerOptions = append(readerOptions, metric.WithReader(promExporter))

	meter := metric.NewMeterProvider(readerOptions...)
	t.shutdownFuncs = append(t.shutdownFuncs, meter.Shutdown)
	return meter
}

func (t *Telemetry) buildResource(ctx context.Context) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("stoke"),
			semconv.ServerAddress(Ctx(ctx).Server.Address),
			semconv.ServerPort(Ctx(ctx).Server.Port),
		),
	)
}
