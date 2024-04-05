package main

import (
	"context"
	"stoke/internal/cfg"
	"stoke/internal/ctx"
	"stoke/internal/ent"
	"stoke/internal/tel"
	"stoke/internal/usr"
	"stoke/internal/web"
)

func main() {
	// TODO command line flags for config
	config := cfg.FromFile("config.yaml") 

	setupLoggers(*config)

	logger.Debug().Interface("config", config).Msg("Config Loaded")
	logger.Info().Msg("Starting Stoke Server...")

	dbClient := createDBClient(*config)
	issuer := createTokenIssuer(*config, dbClient)
	usrProvider := createUserProvider(*config, dbClient)

	globalContext := &ctx.Context{
		Config:       *config,
		Issuer:       issuer,
		DB:           dbClient,
		UserProvider: usrProvider,
		AppContext:   context.Background(),
	}
	
	otel := &tel.OTEL{
		Context: globalContext,
	}

	server := web.Server {
		Context: globalContext,
		OTEL: otel,
	}

	globalContext.OnStartup(otel.Init)
	globalContext.OnStartup(usrProvider.Init)
	globalContext.OnStartup(issuer.Init)
	globalContext.OnStartup(server.Init)

	if err := globalContext.Startup() ; err != nil {
		shutdownErr := globalContext.GracefulShutdown()
		logger.Fatal().Err(err).AnErr("shutdownErr", shutdownErr).Msg("Could not initialize context")
	}

	if err := server.Run(); err != nil {
		logger.Error().Err(err).Msg("Error stopped the server")
	}

	err := globalContext.GracefulShutdown()
	logger.Info().Err(err).Msg("Stoke Server Terminated.")
}



func createUserProvider(_ cfg.Config, dbClient *ent.Client) usr.Provider {
	// TODO need to be able to configure which provider to use
	// Will always have user
	localUserProvider := usr.LocalProvider{
		DB: dbClient,
	}
	err := localUserProvider.Init()
	if err != nil {
		logger.Error().Err(err).Msg("Could not initialize local user database")
	}

	return localUserProvider
}
