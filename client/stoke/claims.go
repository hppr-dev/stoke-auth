package stoke

import (
	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
)

type stokeClaimBuilder struct {
	validator *Claims
}

// Allows for the use of stoke.Token().Requires("myclaim", "myvalue").ForAccess() to validate tokens
func WithToken() stokeClaimBuilder {
	return stokeClaimBuilder{
		validator: &Claims{
			requiredClaims : make(map[string]string),
		},
	}
}

func (b stokeClaimBuilder) Requires(name, value string) stokeClaimBuilder {
	b.validator.requiredClaims[name] = value
	return b
}

func (b stokeClaimBuilder) ForAccess() *Claims {
	return b.validator
}

type Claims struct {
	StokeClaims map[string]string `json:"stk"`
	jwt.RegisteredClaims
	requiredClaims map[string]string
}

func (s Claims) Validate() error {
	for k, v := range s.requiredClaims {
		customVal, ok := s.StokeClaims[k]
		if !ok || customVal != v {
			return errors.New("User lacks the required claim to access resource.")
		}
	}
	return nil
}
