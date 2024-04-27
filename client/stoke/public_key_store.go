package stoke

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// PublicKeyStore are responsible for pulling and storing public keys from a stoke server
type PublicKeyStore interface {
	ParseClaims(context.Context, string, *Claims, ...jwt.ParserOption) (*jwt.Token, error)
}
