package stoke

import (
	"time"

	"github.com/go-faster/jx"
	"github.com/golang-jwt/jwt/v5"
)

// stoke.Claims is used to hold and validate claims on a JWT.
// StokeClaims hold all custom claims that are set.
// It is marshalled using custom logic to place all StokeClaims and RegisteredClaims at the same level in the resulting token.
type Claims struct {
	StokeClaims map[string]string `json:"-"`
	jwt.RegisteredClaims
	requiredClaimPreds []claimPredicate
	alternateClaims []*Claims
}

// claimPredicates are custom functions to verify claims satisfy conditions
type claimPredicate func(c Claims) bool

// Creates an claim object with the same RegisteredClaims and an empty StokeClaims map.
func (c *Claims) New() *Claims {
	return &Claims{
		RegisteredClaims: c.RegisteredClaims,
		StokeClaims: make(map[string]string),
	}
}

// Implements json.Marshaler.
// Custom logic to place StokeClaims and RegisteredClaims at the same level in json.
func (c *Claims) MarshalJSON() ([]byte, error) {
	encoder := jx.GetEncoder()

	if _, ok := c.StokeClaims[""]; ok {
		delete(c.StokeClaims, "")
	}

	encoder.ObjStart()

	iss := c.RegisteredClaims.Issuer
	if iss != "" {
		encoder.FieldStart("iss")
		encoder.Str(iss)
	}

	sub := c.RegisteredClaims.Subject
	if sub != "" {
		encoder.FieldStart("sub")
		encoder.Str(sub)
	}

	aud := c.RegisteredClaims.Audience
	if len(aud) > 0 {
		encoder.FieldStart("aud")
		encoder.ArrStart()
		for _, a := range aud {
			encoder.Str(a)
		}
		encoder.ArrEnd()
	}

	if c.RegisteredClaims.ExpiresAt != nil {
		exp := c.RegisteredClaims.ExpiresAt.Unix()
		if exp > 0 {
			encoder.FieldStart("exp")
			encoder.Int64(exp)
		}
	}

	if c.RegisteredClaims.NotBefore != nil {
		nbf := c.RegisteredClaims.NotBefore.Unix()
		if nbf > 0 {
			encoder.FieldStart("nbf")
			encoder.Int64(nbf)
		}
	}

	if c.RegisteredClaims.IssuedAt != nil {
		iat := c.RegisteredClaims.IssuedAt.Unix()
		if iat > 0 {
			encoder.FieldStart("iat")
			encoder.Int64(iat)
		}
	}

	jti := c.RegisteredClaims.ID
	if jti != "" {
		encoder.FieldStart("jti")
		encoder.Str(jti)
	}

	for name, value := range c.StokeClaims {
		encoder.FieldStart(name)
		encoder.Str(value)
	}
	encoder.ObjEnd()

	return encoder.Bytes(), nil
}

// Implements json.Unmarshaler.
// Custom logic to parse StokeClaims and RegisteredClaims at the same level in json.
func (c *Claims) UnmarshalJSON(b []byte) error {
	decoder := jx.DecodeBytes(b)
	c.StokeClaims = make(map[string]string)

	err := decoder.Obj(func(d *jx.Decoder, key string) (oErr error) {
		var i int64
		switch(key) {
			// Registered Claims
			case "iss":
				c.RegisteredClaims.Issuer, oErr = d.Str()
			case "sub":
				c.RegisteredClaims.Subject, oErr = d.Str()
			case "aud":
				oErr = d.Arr(func(d *jx.Decoder) error {
					s, aErr := d.Str()
					c.RegisteredClaims.Audience = append(c.RegisteredClaims.Audience, s)
					return aErr
				})
			case "exp":
				i, oErr = d.Int64()
				c.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(i, 0 ))
			case "nbf":
				i, oErr = d.Int64()
				c.RegisteredClaims.NotBefore = jwt.NewNumericDate(time.Unix(i, 0))
			case "iat":
				i, oErr = d.Int64()
				c.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Unix(i, 0))
			case "jti":
				c.RegisteredClaims.ID, oErr = d.Str()
			default:
				c.StokeClaims[key], oErr = d.Str()
		}
		return
	})

	return err
}
