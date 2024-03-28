package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"time"

	"stoke/internal/cfg"
	"stoke/internal/ctx"
	"stoke/internal/ent"
	"stoke/internal/key"
	"stoke/internal/usr"
	"stoke/internal/web"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting Stoke Server...")

	config := cfg.Config{
		Server: cfg.Server {
			Address : "",
			Port: 8080,
		},
	}

	dbClient := connectToDB(config)

	localUserProvider := usr.LocalProvider{
		DB: dbClient,
	}
	err := localUserProvider.Init()
	if err != nil {
		log.Panicf("Could not initialize local user database: %v", err)
	}

	appContext := &ctx.Context{
		Config:       config,
		Issuer:       createTokenIssuer(config, dbClient),
		DB:           dbClient,
		UserProvider: localUserProvider,
	}

	server := web.Server {
		Context: appContext,
	}

	server.Init()
	if err := server.Run(); err != nil {
		log.Printf("An error stopped the server: %v", err)
	}
	
	log.Println("Stoke Server Terminated.")
}

func connectToDB(_ cfg.Config) *ent.Client {
	// Type and connection string should be configurable
	dbClient, err := ent.Open("sqlite3", "file:stoke.db?cache=shared&_fk=1")
	if err != nil {
		log.Panicf("Could not connect to database: %v", err)
	}

	dbClient.Schema.Create(context.Background())
	return dbClient
}

func createTokenIssuer(_ cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	// Also should be an option to persist or not perist in db
	// type, key and token duration must be configurable
	cache := key.KeyCache[*ecdsa.PrivateKey]{
		KeyDuration:   time.Hour * 3,
		TokenDuration: time.Minute * 30,
	}
	// This will depend on the configuration of certificate type
	pair := &key.ECDSAKeyPair{
		NumBits : 256,
	}

	err := cache.Bootstrap(dbClient, pair)
	if err != nil {
		log.Panicf("Could not bootstrap key cache! %v", err)
	}

	return &key.AsymetricTokenIssuer[*ecdsa.PrivateKey]{
		KeyCache: cache,
	}
}
