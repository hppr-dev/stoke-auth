package stoke

import (
)

type PublicKeyStore interface {
	Init() error
	ValidateClaims(string, *ClaimsValidator) bool
}

func DefaultPublicKeyStore(endpoint string) PublicKeyStore {
	return &WebPublicKeyStore{
		Endpoint: endpoint,
	}
}
