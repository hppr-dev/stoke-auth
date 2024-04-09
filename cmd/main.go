package main

import (
	"context"
	"errors"
	"stoke/internal/cfg"
	"stoke/internal/web"

	"github.com/rs/zerolog"
)

func main() {
	// TODO command line flags for config
	config := cfg.FromFile("config.yaml") 

	rootCtx := config.WithContext(context.Background())

	logger := zerolog.Ctx(rootCtx)

	logger.Debug().Interface("config", config).Msg("Config Loaded")
	logger.Info().
		Str("addr",config.Server.Address).
		Int("port", config.Server.Port).
		Str("privateKey", config.Server.TLSPrivateKey).
		Str("publicCert", config.Server.TLSPublicCert).
		Msg("Starting Stoke Server...")

	server := web.NewServer(rootCtx)

	shutdownFuncs, err := config.Telemetry.Initialize(rootCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not initialize telemetry")
	}

	if server.TLSConfig != nil {
		if err := server.ListenAndServeTLS("","") ; err != nil {
			logger.Error().Err(err).Msg("An error occurred with the TLS server")
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			logger.Error().Err(err).Msg("An error occurred with the http server")
		}
	}

	err = nil
	for _, f := range shutdownFuncs {
		errors.Join(err, f(rootCtx))
	}

	logger.Info().Err(err).Msg("Stoke Server Terminated.")
}

