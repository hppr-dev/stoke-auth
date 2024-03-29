package stoke

import "github.com/golang-jwt/jwt/v5"

type PublicKeyStore interface {
	Init() error
	ValidateClaims(string, *ClaimsValidator, ...jwt.ParserOption) bool
}

func DefaultPublicKeyStore(endpoint string) PublicKeyStore {
	return &WebPublicKeyStore{
		Endpoint: endpoint,
	}
}
