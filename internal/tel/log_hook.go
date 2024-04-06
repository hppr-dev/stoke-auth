package tel

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)


type LogHook struct{
	Ctx context.Context
}

func (l LogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level >= zerolog.GlobalLevel() {
		span := trace.SpanFromContext(l.Ctx)
		span.AddEvent(fmt.Sprintf("LOG %s: %s", level.String(), msg))
	}
}
