package main

import (
	"stoke/internal/cfg"
	"stoke/internal/ctx"
	"stoke/internal/ent"
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

	appContext := &ctx.Context{
		Config:       *config,
		Issuer:       createTokenIssuer(*config, dbClient),
		DB:           dbClient,
		UserProvider: createUserProvider(*config, dbClient),
	}

	if err := appContext.Issuer.Init() ; err != nil {
		logger.Fatal().Err(err).Msg("Could not initialize token issuer")
	}

	server := web.Server {
		Context: appContext,
	}

	server.Init()
	if err := server.Run(); err != nil {
		logger.Error().Err(err).Msg("An error stopped the server")
	}
	
	logger.Info().Msg("Stoke Server Terminated.")
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
