package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"stoke/internal/cfg"
	_ "stoke/internal/ent/runtime"
	"stoke/internal/usr"
	"stoke/internal/web"

	"github.com/rs/zerolog"
)

func main() {
	flagSet := flag.NewFlagSet("flags", flag.ExitOnError)
	dbInitFile := flagSet.String("dbinit", "", "Database initialization file (overrides what is in config.yaml)")
	configFile := flagSet.String("config", "config.yaml", "Configuration file to use")

	var allFlags []string
	validateOnly := false
	migrateOnly := false
	hashOnly := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "hash-password":
			hashOnly = true
		case "validate":
			validateOnly = true
		case "migrate":
			migrateOnly = true
		default:
			allFlags = append(allFlags, arg)
		}
	}

	if hashOnly {
		getAndHashPassword()
		return
	}

	if err := flagSet.Parse(allFlags); err != nil {
		fmt.Printf("Failed to parse flags: %v\n", err)
		os.Exit(1)
	}

	config := cfg.FromFile(*configFile) 

	if validateOnly {
		fmt.Printf("Config Validated: %+v", config)
		return
	}

	if *dbInitFile != "" {
		config.Users.UserInitFile = *dbInitFile
	}

	rootCtx := config.WithContext(context.Background())
	logger := zerolog.Ctx(rootCtx)

	logger.Debug().
		Str("configFile", *configFile).
		Str("dbInitFile", *dbInitFile).
		Interface("config", config).
		Msg("Config Loaded")


	if err := usr.ProviderFromCtx(rootCtx).CheckCreateForSuperUser(rootCtx); err != nil {
		logger.Error().
			Err(err).
			Msg("Could not check/create super user")
	}

	if migrateOnly {
		return
	}

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
		err = errors.Join(err, f(rootCtx))
	}

	logger.Info().Err(err).Msg("Stoke Server Terminated.")
}

func getAndHashPassword() {
	var pass string
	fmt.Println("Creating password hash for db-init file...")
	fmt.Print("password:\033[8m")
	if _, err := fmt.Scanln(&pass) ; err != nil {
		fmt.Printf("Could not read password: %v", err)
		os.Exit(1)
	}
	salt := usr.GenSalt()
	hash := usr.HashPass(pass, salt)
	fmt.Println("\033[28mAdd the following to the db-init yaml file:")
	fmt.Printf("\npassword_hash: %s\npassword_salt: %s\n", hash, salt)
}
