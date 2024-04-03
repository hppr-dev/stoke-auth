package stoke

import "github.com/golang-jwt/jwt/v5"

type PublicKeyStore interface {
	Init() error
	ParseClaims(string, *Claims, ...jwt.ParserOption) (*jwt.Token, error)
}

func DefaultPublicKeyStore(endpoint string) PublicKeyStore {
	return &WebPublicKeyStore{
		Endpoint: endpoint,
	}
}
