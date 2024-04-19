package stoke

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type PublicKeyStore interface {
	Init(context.Context) error
	ParseClaims(context.Context, string, *Claims, ...jwt.ParserOption) (*jwt.Token, error)
}

func DefaultPublicKeyStore(endpoint string) PublicKeyStore {
	return &WebCachePublicKeyStore{
		Endpoint: endpoint,
	}
}
