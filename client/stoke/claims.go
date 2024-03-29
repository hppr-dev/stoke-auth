package stoke

import (
	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
)

type stokeClaimBuilder struct {
	validator *ClaimsValidator
}

// Allows for the use of stoke.Claims().Require("myclaim", "myvalue").Validator() to validate tokens
func Claims() stokeClaimBuilder {
	return stokeClaimBuilder{
		validator: &ClaimsValidator{
			requiredClaims : make(map[string]string),
		},
	}
}

func (b stokeClaimBuilder) Require(name, value string) stokeClaimBuilder {
	b.validator.requiredClaims[name] = value
	return b
}

func (b stokeClaimBuilder) Validator() *ClaimsValidator {
	return b.validator
}

type ClaimsValidator struct {
	CustomClaims map[string]string `json:"stk"`
	jwt.RegisteredClaims
	requiredClaims map[string]string
}

func (s ClaimsValidator) Validate() error {
	for k, v := range s.requiredClaims {
		customVal, ok := s.CustomClaims[k]
		if !ok || customVal != v {
			return errors.New("User lacks the required claim to access resource.")
		}
	}
	return nil
}
