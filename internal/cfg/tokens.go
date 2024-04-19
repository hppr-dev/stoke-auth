package cfg

import (
	"context"
	"regexp"
	"stoke/internal/key"
	"strconv"
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
	t.TokenDuration = strToDuration(t.TokenDurationStr)
	t.KeyDuration = strToDuration(t.KeyDurationStr)
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

	return context.WithValue(ctx, "issuer", issuer)
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
	cache := key.KeyCache[P]{
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

var durationRegex *regexp.Regexp = regexp.MustCompile(`(\d+)([sSmMhHdDyY])`)

func strToDuration(s string) time.Duration {
	matches := durationRegex.FindStringSubmatch(s)
	if len(matches) != 3 {
		panic("Duration string did not match regex [0-9]+[sSmMhHdDyY")
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
