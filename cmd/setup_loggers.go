package main

import (
	"fmt"
	"io"
	"os"

	"stoke/internal/cfg"
	"stoke/internal/key"
	"stoke/internal/usr"
	"stoke/internal/web"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

func setupLoggers(conf cfg.Config) {
	logger = log.Logger
	var rootWriter io.Writer
	var writers []io.Writer = []io.Writer{ os.Stdout }

	if conf.Logging.LogFile != "" {
		f, err := os.Create(conf.Logging.LogFile)
		if err != nil {
			panic(fmt.Sprintf("Could not open log file: %s, %v", conf.Logging.LogFile, err))
		}
		writers = append(writers, f)
	}

	rootWriter = io.MultiWriter(writers...)
	if conf.Logging.PrettyConsole {
		rootWriter = zerolog.ConsoleWriter{ Out: rootWriter }
	} 
	logger = log.Output(rootWriter)

	switch conf.Logging.Level {
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

	usr.SetLogger(logger)
	web.SetLogger(logger)
	key.SetLogger(logger)
}
