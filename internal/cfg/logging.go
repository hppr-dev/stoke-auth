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
	// File to write logs to
	LogFile       string `json:"log_file"`
	// Whether to prettify the output. Set to false to log JSON.
	PrettyConsole bool   `json:"pretty_console"`
}

func (l Logging) withContext(ctx context.Context) context.Context {
	logger := log.Logger

	var rootWriter io.Writer
	var writers []io.Writer = []io.Writer{ os.Stdout }

	if l.LogFile != "" {
		f, err := os.Create(l.LogFile)
		if err != nil {
			panic(fmt.Sprintf("Could not open log file: %s, %v", l.LogFile, err))
		}
		writers = append(writers, f)
	}

	rootWriter = io.MultiWriter(writers...)
	if l.PrettyConsole {
		rootWriter = zerolog.ConsoleWriter{ Out: rootWriter }
	} 
	logger = log.Output(rootWriter)

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
