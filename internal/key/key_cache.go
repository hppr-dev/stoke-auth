package key

import (
	"context"
	"encoding/json"
	"stoke/client/stoke"
	"stoke/internal/ent"
	"stoke/internal/ent/privatekey"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type KeyCache[P PrivateKey] struct {
	keys []KeyPair[P]
	KeyDuration time.Duration
	TokenDuration time.Duration
}

// Implements stoke.PublicKeyStore
func (c *KeyCache[P]) Init() error {
	go c.goManage()
	return nil
}

func (c *KeyCache[P]) goManage() {
	// TODO Manage keys (rotate, clean, etc)
	logger.Info().Msg("Starting key cache management...")
}

func (c *KeyCache[P]) CurrentKey() KeyPair[P] {
	return c.keys[len(c.keys) - 1]
}

type publicJson struct {
	Text    string    `json:"text"`
	Expires int64     `json:"expires"`
	Renews  int64     `json:"renews"`
	Method  string    `json:"method"`
}

func (c *KeyCache[P]) PublicKeys() ([]byte, error) {
	out := make([]publicJson, len(c.keys))
	for i, k := range c.keys {
		out[i] = publicJson{
			Text:    k.PublicString(),
			Expires: k.ExpiresAt().Unix(),
			Renews:  k.RenewsAt().Unix(),
			Method:  k.SigningMethod().Alg(),
		}
	}
	return json.Marshal(out)
}

func (c *KeyCache[P]) Generate() error {
	logger.Debug().Msg("Generating new key...")

	newKey := new(KeyPair[P])
	err := (*newKey).Generate()
	if err != nil {
		logger.Error().Err(err).Msg("Could not generate key")
		return err
	}

	expires := time.Now().Add(c.KeyDuration)
	renews  := time.Now().Add(c.KeyDuration).Add(-c.TokenDuration)

	logger.Debug().
		Time("expires", expires).
		Time("renews", renews).
		Msg("Generated new key.")

	(*newKey).SetExpires(expires)
	(*newKey).SetRenews(renews)

	c.keys = append(c.keys, *newKey)

	return nil
}

func (c *KeyCache[P]) Bootstrap(db *ent.Client, pair KeyPair[P]) error {
	logger.Info().Msg("Bootstraping key cache.")
	now := time.Now()
	pk, err := db.PrivateKey.Query().
		Order(privatekey.ByExpires(sql.OrderDesc())).
		First(context.Background())

	if err != nil || pk.Expires.Before(now) {
		logger.Info().Msg("Could not retrieve private key. Generating a new one.")
		pair.Generate()

		pk, err = db.PrivateKey.Create().
			SetText(pair.Encode()).
			SetExpires(now.Add(c.KeyDuration)).
			SetRenews(now.Add(c.KeyDuration).Add(-c.TokenDuration)).
			Save(context.Background())
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
	pair.SetRenews(pk.Renews)

	c.keys = append(c.keys, pair)
	return nil
}

func (c *KeyCache[P]) Clean() {
	logger.Info().Msg("Cleaning key cache...")

	now := time.Now()
	var valid []KeyPair[P]
	for _, e := range c.keys {
		if e.ExpiresAt().Before(now) {
			valid = append(valid, e)
		}
	}
	c.keys = valid

	logger.Debug().Int("remainingKeys", len(c.keys)).Msg("Done Cleaning.")
}

// Implements stoke.PublicKeyStore
func (c *KeyCache[P]) ValidateClaims(token string, claims *stoke.ClaimsValidator, parserOpts ...jwt.ParserOption) bool {
	jwtToken, err := jwt.ParseWithClaims(token, claims, c.publicKeys, parserOpts...)
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to validate claims")
		return false
	}

	logger.Debug().
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
	return jwtToken.Valid
}

func (c *KeyCache[P]) publicKeys(_ *jwt.Token) (interface{}, error) {
	pkeys := jwt.VerificationKeySet{}
	for _, p := range c.keys {
		pkeys.Keys = append(pkeys.Keys, p.PublicKey())
	}
	return pkeys, nil
}
