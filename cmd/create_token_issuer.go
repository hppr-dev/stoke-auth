package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/key"
)

func createTokenIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	// TODO Also should be an option to persist or not perist in db
	// type, key and token duration must be configurable
	switch conf.Tokens.Algorithm {
	case "ECDSA", "ecdsa":
		return createECDSAIssuer(conf, dbClient)
	case "EdDSA", "eddsa":
		return createEdDSAIssuer(conf, dbClient)
	case "RSA", "rsa":
		return createRSAIssuer(conf, dbClient)
	}
	logger.Fatal().Str("algorithm", conf.Tokens.Algorithm).Msg("Unsupported algorithm")
	return nil
}

func createECDSAIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	cache := key.KeyCache[*ecdsa.PrivateKey]{
		KeyDuration:   conf.Tokens.KeyDuration,
		TokenDuration: conf.Tokens.TokenDuration,
	}
	pair := &key.ECDSAKeyPair{
		NumBits : conf.Tokens.NumBits,
	}

	err := cache.Bootstrap(dbClient, pair)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap key cache")
	}

	return &key.AsymetricTokenIssuer[*ecdsa.PrivateKey]{
		KeyCache: cache,
	}
}

func createEdDSAIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	cache := key.KeyCache[ed25519.PrivateKey]{
		KeyDuration:   conf.Tokens.KeyDuration,
		TokenDuration: conf.Tokens.TokenDuration,
	}
	pair := &key.EdDSAKeyPair{}

	err := cache.Bootstrap(dbClient, pair)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap key cache")
	}

	return &key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: cache,
	}
}

func createRSAIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	cache := key.KeyCache[*rsa.PrivateKey]{
		KeyDuration:   conf.Tokens.KeyDuration,
		TokenDuration: conf.Tokens.TokenDuration,
	}
	pair := &key.RSAKeyPair{}

	err := cache.Bootstrap(dbClient, pair)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap key cache")
	}

	return &key.AsymetricTokenIssuer[*rsa.PrivateKey]{
		KeyCache: cache,
	}
}


