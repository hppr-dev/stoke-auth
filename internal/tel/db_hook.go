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


func DatabaseMutationTelemetry(ctx context.Context) func(ent.Mutator) ent.Mutator {
	logger := zerolog.Ctx(ctx)

	return func(next ent.Mutator) ent.Mutator{
		timerHistogram, err := GetMeter().Int64Histogram(
			"stoke_database_time_histogram",
			metric.WithDescription("Histogram of database operation time"),
		)
		if err != nil {
			logger.Error().Err(err).Msg("Could not create database mutation metrics instrumentation")
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
		t.timerHistogram.Record(ctx, milliDiff,
			metric.WithAttributes(
				attribute.String("operation", "mutate"),
				attribute.String("unit", "milli"),
			),
		) 
		span.End()
	}()
	return t.next.Mutate(ctx, m)
}

func DatabaseReadTelemetry(ctx context.Context) ent.InterceptFunc {
	logger := zerolog.Ctx(ctx)

	return func(next ent.Querier) ent.Querier {
		timerHistogram, err := GetMeter().Int64Histogram(
			"stoke_database_time_histogram",
			metric.WithDescription("Histogram of database operation time"),
		)
		if err != nil {
			logger.Error().Err(err).Msg("Could not create database read metrics instrumentation")
			return next
		}

		return telemetryQuerier{
			timerHistogram: timerHistogram,
			next: next,
			tracer: GetTracer(),
		}
	}
}

type telemetryQuerier struct {
	tracer trace.Tracer
	timerHistogram metric.Int64Histogram
	next ent.Querier
}

func (t telemetryQuerier) Query(ctx context.Context, q ent.Query) (ent.Value, error) {
	tableName := "unknown"
	switch q.(type) {
	case *ent.UserQuery:
		tableName = "users"
	case *ent.ClaimGroupQuery:
		tableName = "claim_group"
	case *ent.ClaimQuery:
		tableName = "claim"
	case *ent.GroupLinkQuery:
		tableName = "group_link"
	}

	ctx, span := t.tracer.Start(ctx, "Database.Read",
		trace.WithAttributes(
			attribute.String("tableName", tableName),
		),
	)

	start := time.Now()
	defer func() {
		uDiff := time.Now().UnixMicro() - start.UnixMicro()
		t.timerHistogram.Record(ctx, uDiff,
			metric.WithAttributes(
				attribute.String("operation", "read"),
				attribute.String("unit", "micro"),
			),
		) 
		span.End()
	}()
	return t.next.Query(ctx, q)
}
