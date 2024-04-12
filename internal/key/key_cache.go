package key

import (
	"context"
	"encoding/json"
	"fmt"
	"stoke/client/stoke"
	"stoke/internal/ent"
	"stoke/internal/ent/privatekey"
	"stoke/internal/tel"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

type KeyCache[P PrivateKey] struct {
	activeKey int
	keys []KeyPair[P]
	Ctx context.Context
	KeyDuration time.Duration
	TokenDuration time.Duration
}

// Implements stoke.PublicKeyStore
func (c *KeyCache[P]) Init(ctx context.Context) error {
	c.activeKey = 0
	go c.goManage(ctx)
	return nil
}

func (c *KeyCache[P]) goManage(ctx context.Context) {
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("Starting key cache management...")
	for {
		nextExpire := c.CurrentKey().ExpiresAt().Sub(time.Now())
		nextRenew  := nextExpire - (c.TokenDuration * 2)

		time.Sleep(nextRenew)
		sCtx, span := tel.GetTracer().Start(ctx, "KeyCache.Rotation")
		sLogger := logger.With().
			Str("component", "KeyCache.Management").
			Logger().
			Hook(tel.LogHook{ Ctx : sCtx } )
		sCtx = sLogger.WithContext(sCtx)

		c.Generate(sCtx)

		time.Sleep(c.TokenDuration)
		sLogger.Info().Msg("Activating new key...")
		c.activeKey += 1

		time.Sleep(c.TokenDuration)
		c.Clean(sCtx)

		span.End()
	}
}

func (c *KeyCache[P]) CurrentKey() KeyPair[P] {
	return c.keys[c.activeKey]
}

func (c *KeyCache[P]) PublicKeys(ctx context.Context) ([]byte, error) {
	ctx, span := tel.GetTracer().Start(ctx, "KeyCache.PublicKeys")
	defer span.End()

	now := time.Now()
	jwks := make([]*stoke.JWK, len(c.keys))
	for i, k := range c.keys {
		jwks[i] = stoke.CreateJWK().FromPublicKey(k.PublicKey())
		jwks[i].KeyId = fmt.Sprintf("p-%d", i)
	}
	expireTime := c.CurrentKey().ExpiresAt()
	clientPullTime := expireTime.Add( ( c.TokenDuration * -3 ) / 2)
	if now.After(clientPullTime) {
		// Clients should refresh after the current key expires
		clientPullTime = expireTime.Add(100 * time.Millisecond)
	}

	return json.Marshal(stoke.JWKSet{
		Expires: clientPullTime,
		Keys: jwks,
	})
}

func (c *KeyCache[P]) Generate(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "KeyCache.Generate")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Msg("Generating new key...")

	if len(c.keys) == 0 {
		logger.Fatal().Msg("Unable to generate keys. No keys in keystore!")
	}

	newKey, err := c.keys[0].Generate()
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not generate key")
		return err
	}

	expires := time.Now().Add(c.KeyDuration)
	newKey.SetExpires(expires)
	c.keys = append(c.keys, newKey)

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Int("numKeys", len(c.keys)).
		Time("expires", expires).
		Dur("keyDuration", c.KeyDuration).
		Dur("tokenDuration", c.TokenDuration).
		Str("publicKey", newKey.PublicString()).
		Msg("Generated new key.")

	return nil
}

func (c *KeyCache[P]) Bootstrap(ctx context.Context, pair KeyPair[P]) error {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "KeyCache.Bootstrap")
	defer span.End()

	logger.Info().
		Msg("Bootstraping key cache.")

	var err error

	db := ent.FromContext(c.Ctx)
	now := time.Now()

	pk, err := db.PrivateKey.Query().
		Order(privatekey.ByExpires(sql.OrderDesc())).
		First(c.Ctx)

	if err != nil || pk.Expires.Before(now) {
		logger.Info().
			Msg("Could not retrieve private key. Generating a new one.")

		pair, err = pair.Generate()
		if err != nil {
			logger.Error().Err(err).Msg("Could not generate private key")
			return err
		}

		pk, err = db.PrivateKey.Create().
			SetText(pair.Encode()).
			SetExpires(now.Add(c.KeyDuration)).
			Save(c.Ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Could not save private key")
			return err
		}
	} else {
		err := pair.Decode(pk.Text)
		if err != nil {
			logger.Error().Err(err).Msg("Could not decode private key text from database")
			return err
		}
	}

	pair.SetExpires(pk.Expires)

	c.keys = append(c.keys, pair)
	return nil
}

func (c *KeyCache[P]) Clean(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "KeyCache.Clean")
	defer span.End()

	logger.Info().
		Func(otelzerolog.AddTracingContext(span)).
		Msg("Cleaning key cache...")
	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Func(func(e *zerolog.Event) {
			pkeyStrs := make([]string, len(c.keys))
			for i, k := range c.keys {
				pkeyStrs[i] = k.PublicString()
			}
			e.Strs("publicKeys", pkeyStrs)
		}).
		Msg("Starting clean")

	now := time.Now()
	var valid []KeyPair[P]
	for _, e := range c.keys {
		if e.ExpiresAt().After(now) {
			valid = append(valid, e)
		}
	}
	c.keys = valid
	c.activeKey = len(c.keys) - 1

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Func(func(e *zerolog.Event) {
			pkeyStrs := make([]string, len(c.keys))
			for i, k := range c.keys {
				pkeyStrs[i] = k.PublicString()
			}
			e.Strs("publicKeys", pkeyStrs)
		}).Msg("Finished cleaning.")
}

// Implements stoke.PublicKeyStore
func (c *KeyCache[P]) ParseClaims(ctx context.Context, token string, claims *stoke.Claims, parserOpts ...jwt.ParserOption) (*jwt.Token, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "KeyCache.ParseClaims")
	defer span.End()

	jwtToken, err := jwt.ParseWithClaims(token, claims, c.publicKeys, parserOpts...)
	if err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Str("token", token).
			Err(err).
			Msg("Failed to validate claims")
		return nil, err
	}

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Bool("valid", jwtToken.Valid).
		Str("alg", jwtToken.Method.Alg()).
		Func(func(e *zerolog.Event) {
			issued, _ := jwtToken.Claims.GetIssuedAt()
			if issued != nil {
				e.Time("issued", issued.Time)
			}

			expires, _ := jwtToken.Claims.GetExpirationTime()
			if expires != nil {
				e.Time("expires", expires.Time)
			}
		}).
		Msg("Parsed Token")
	return jwtToken, err
}

func (c *KeyCache[P]) publicKeys(_ *jwt.Token) (interface{}, error) {
	pkeys := jwt.VerificationKeySet{}
	for _, p := range c.keys {
		pkeys.Keys = append(pkeys.Keys, p.PublicKey())
	}
	return pkeys, nil
}
