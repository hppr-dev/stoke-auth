package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"regexp"
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/key"
	"strconv"
	"time"
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
	panic("Unrecoverable")
	return nil
}

func createECDSAIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	cache := key.KeyCache[*ecdsa.PrivateKey]{
		KeyDuration:   strToDuration(conf.Tokens.KeyDuration),
		TokenDuration: strToDuration(conf.Tokens.TokenDuration),
	}
	pair := &key.ECDSAKeyPair{
		NumBits : conf.Tokens.NumBits,
	}

	err := cache.Bootstrap(dbClient, pair)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap key cache")
		panic("Unrecoverable")
	}

	return &key.AsymetricTokenIssuer[*ecdsa.PrivateKey]{
		KeyCache: cache,
	}
}

func createEdDSAIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	cache := key.KeyCache[ed25519.PrivateKey]{
		KeyDuration:   strToDuration(conf.Tokens.KeyDuration),
		TokenDuration: strToDuration(conf.Tokens.TokenDuration),
	}
	pair := &key.EdDSAKeyPair{}

	err := cache.Bootstrap(dbClient, pair)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap key cache")
		panic("Unrecoverable")
	}

	return &key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: cache,
	}
}

func createRSAIssuer(conf cfg.Config, dbClient *ent.Client) key.TokenIssuer {
	cache := key.KeyCache[*rsa.PrivateKey]{
		KeyDuration:   strToDuration(conf.Tokens.KeyDuration),
		TokenDuration: strToDuration(conf.Tokens.TokenDuration),
	}
	pair := &key.RSAKeyPair{}

	err := cache.Bootstrap(dbClient, pair)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap key cache")
		panic("Unrecoverable")
	}

	return &key.AsymetricTokenIssuer[*rsa.PrivateKey]{
		KeyCache: cache,
	}
}

var durationRegex *regexp.Regexp = regexp.MustCompile(`(\d+)([sSmMhHdDyY])`)

func strToDuration(s string) time.Duration {
	matches := durationRegex.FindStringSubmatch(s)
	if len(matches) != 3 {
		logger.Fatal().Str("durationString", s).Msg("Was not parsable")
		panic("Unrecoverable")
	}
	num, _ := strconv.Atoi(matches[1])
	dur := time.Duration(num)
	switch matches[2] {
	case "s", "S":
		return time.Second * dur
	case "m", "M":
		return time.Minute * dur
	case "h", "H":
		return time.Hour * dur
	case "d", "D":
		return time.Hour * 24 * dur
	case "y", "Y":
		return time.Hour * 24 * 265 * dur
	}
	panic("Unreachable. If it reaches here, it means that durationRegex is broken.")
}
