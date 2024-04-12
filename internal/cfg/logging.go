package cfg

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logging struct {
	// ONE of trace, debug, info, warn, error, fatal, panic
	Level         string `json:"level"`
	// File to write logs to. Leave empty to not log to file
	LogFile       string `json:"log_file"`
	// Whether to prettify the output. Set to false to log JSON.
	PrettyStdout bool   `json:"pretty_stdout"`
	// Whether to write logs to stdout.
	WriteToStdout bool   `json:"write_to_stdout"`
}

func (l Logging) withContext(ctx context.Context) context.Context {
	logger := log.Logger

	var writers []io.Writer = []io.Writer{ }

	if l.WriteToStdout {
		if l.PrettyStdout {
			writers = append(writers, zerolog.ConsoleWriter{ Out: os.Stdout } )
		} else {
			writers = append(writers, os.Stdout )
		}
	}

	if l.LogFile != "" {
		f, err := os.OpenFile(l.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(fmt.Sprintf("Could not open log file: %s, %v", l.LogFile, err))
		}
		writers = append(writers, f)
	}

	logger = log.Output(io.MultiWriter(writers...))

	switch l.Level {
	case "TRACE", "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		logger = logger.With().Caller().Logger()
	case "DEBUG", "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.With().Caller().Logger()
	case "INFO", "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN", "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL", "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC", "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger.Info().Msg("Log level not specified. Setting to INFO by default.")
	}

	return logger.WithContext(ctx)
}
