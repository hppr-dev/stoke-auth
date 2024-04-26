package cfg

import (
	"context"
	"fmt"
	"stoke/internal/key"
	"time"

	"github.com/rs/zerolog"
)


type Tokens struct {
	// One of RSA, ECDSA, or EdDSA
	Algorithm        string `json:"algorithm"`
	// Only applies for RSA and ECDSA
	NumBits          int    `json:"num_bits"`
	// Whether or not to save the private keys in the database
	PersistKeys      bool   `json:"persist_keys"`
	// How long to keep signing keys alive
	KeyDurationStr   string `json:"key_duration"`
	// How long to issue tokens for
	TokenDurationStr string `json:"token_duration"`
	// Include Not Before Header
	IncludeNotBefore bool     `json:"include_not_before"`
	// Include Issued At Header
	IncludeIssuedAt bool     ` json:"include_issued_at"`
	// Issuer to set on all tokens
	Issuer           string   `json:"issuer"`
	// Subject to set on all tokens
	Subject          string   `json:"subject"`
	// Audience to set on all tokens
	Audience         []string `json:"audience"`
	// List of user identifiers to include in Tokens
	// May have any or all of the following keys: username, first_name, last_name, full_name, email
	UserInfo         map[string]string `json:"user_info"`

	// Maximum number of refreshes per token. Set to 0 for unlimited
	TokenRefreshLimit     int `json:"token_refresh_limit"`
	// Key to hold the token's refresh count. Omit to not include in tokens. Defaults to using the registered jwt id header.
	TokenRefreshCountKey  string `json:"token_refresh_count_key"`

	// Non-parsed fields
	TokenDuration time.Duration `json:"-"`
	KeyDuration time.Duration   `json:"-"`
}

func (t *Tokens) ParseDurations() {
	var err error
	t.TokenDuration, err = time.ParseDuration(t.TokenDurationStr)
	if err != nil {
		panic(fmt.Sprintf("Could not parse duration \"%s\": %v", t.TokenDurationStr, err))
	}
	t.KeyDuration, err = time.ParseDuration(t.KeyDurationStr)
	if err != nil {
		panic(fmt.Sprintf("Could not parse duration \"%s\": %v", t.KeyDurationStr, err))
	}
}

func (t *Tokens) withContext(ctx context.Context) context.Context {
	var issuer key.TokenIssuer

	t.ParseDurations()

	switch t.Algorithm {
	case "ECDSA", "ecdsa":
		issuer = t.createECDSAIssuer(ctx)
	case "EdDSA", "eddsa":
		issuer = t.createEdDSAIssuer(ctx)
	case "RSA", "rsa":
		issuer = t.createRSAIssuer(ctx)
	}

	if issuer == nil {
		zerolog.Ctx(ctx).Fatal().
			Str("component", "cfg.Tokens").
			Str("algorithm", t.Algorithm).
			Msg("Unsupported algorithm")
	}

	if err := issuer.Init(ctx); err != nil {
		zerolog.Ctx(ctx).Fatal().
			Str("component", "cfg.Tokens").
			Err(err).
			Msg("Could not initialize issuer")
	}

	return issuer.WithContext(ctx)
}

func (t *Tokens) createECDSAIssuer(ctx context.Context) key.TokenIssuer {
	return createAsymetricIssuer(t, ctx,
		&key.ECDSAKeyPair{
			NumBits: t.NumBits,
			Logger: zerolog.Ctx(ctx).With().Str("component", "ECDSAKeyPair").Logger(),
		},
	)
}

func (t *Tokens) createEdDSAIssuer(ctx context.Context) key.TokenIssuer {
	return createAsymetricIssuer(t, ctx,
		&key.EdDSAKeyPair{
			Logger: zerolog.Ctx(ctx).With().Str("component", "EdDSAKeyPair").Logger(),
		},
	)
}

func (t *Tokens) createRSAIssuer(ctx context.Context) key.TokenIssuer {
	return createAsymetricIssuer(t, ctx,
		&key.RSAKeyPair{
			NumBits: t.NumBits,
			Logger: zerolog.Ctx(ctx).With().Str("component", "RSAKeyPair").Logger(),
		},
	)
}

func createAsymetricIssuer[P key.PrivateKey](t *Tokens, ctx context.Context, pair key.KeyPair[P]) *key.AsymetricTokenIssuer[P] {
	cache := key.PrivateKeyCache[P]{
		Ctx: augmentContext(ctx, "KeyCache"),
		KeyDuration:   t.KeyDuration,
		TokenDuration: t.TokenDuration,
		PersistKeys:   t.PersistKeys,
	}

	err := cache.Bootstrap(ctx, pair)
	if err != nil {
		zerolog.Ctx(ctx).Fatal().
			Str("component", "cfg.Tokens").
			Err(err).
			Msg("Could not bootstrap key cache")
	}

	return &key.AsymetricTokenIssuer[P]{
		KeyCache: &cache,
		TokenRefreshLimit: t.TokenRefreshLimit,
		TokenRefreshCountKey: t.TokenRefreshCountKey,
	}
}
