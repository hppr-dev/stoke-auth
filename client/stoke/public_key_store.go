package stoke

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// PublicKeyStore are responsible for pulling and storing public keys from a stoke server
type PublicKeyStore interface {
	Init(context.Context) error
	ParseClaims(context.Context, string, *Claims, ...jwt.ParserOption) (*jwt.Token, error)
}

// The DefaultPublicKeyStore is a WebCachePublicKeyStore
func DefaultPublicKeyStore(endpoint string) PublicKeyStore {
	return &WebCachePublicKeyStore{
		Endpoint: endpoint,
	}
}
