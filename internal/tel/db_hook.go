package tel

import (
	"context"
	"stoke/internal/ent"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)


func DatabaseTelemetryMiddleware(ctx context.Context) func(ent.Mutator) ent.Mutator {
	logger := zerolog.Ctx(ctx)

	return func(next ent.Mutator) ent.Mutator{
		timerHistogram, err := GetMeter().Int64Histogram(
			"stoke_database_mutation_millis_histogram",
			metric.WithDescription("Histogram of database operation time in milliseconds"),
			metric.WithUnit("millisecond"),
		)
		if err != nil {
			logger.Error().Err(err).Msg("Could not create metrics instrumentation")
			return next
		}

		return telemetryMutator{
			timerHistogram: timerHistogram,
			next: next,
			tracer: GetTracer(),
		}
	}
}

type telemetryMutator struct {
	tracer trace.Tracer
	timerHistogram metric.Int64Histogram
	next ent.Mutator
}

func (t telemetryMutator) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	ctx, span := t.tracer.Start(ctx, "Database.Mutate",
		trace.WithAttributes(
			attribute.String("schemaType", m.Type()),
			attribute.StringSlice("clearedFields", m.ClearedFields()),
			attribute.StringSlice("addedFields", m.AddedFields()),
			attribute.StringSlice("clearedEdges", m.ClearedEdges()),
			attribute.StringSlice("addedEdges", m.AddedEdges()),
		),
	)
	start := time.Now()
	defer func() {
		milliDiff := time.Now().UnixMilli() - start.UnixMilli()
		t.timerHistogram.Record(ctx, milliDiff) 
		span.End()
	}()
	return t.next.Mutate(ctx, m)
}
