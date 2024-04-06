package tel

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)


type LogHook struct{
	ctx context.Context
}

func (l LogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level >= zerolog.GlobalLevel() {
		span := trace.SpanFromContext(l.ctx)
		span.AddEvent("LOG: " + msg)
	}
}
