package stoke

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PerRequestPublicKeyStore checks the expire time on every request.
// If the expire time is in the past, it refreshes the list of public keys before verifying the token
type PerRequestPublicKeyStore struct {
	BasePublicKeyStore
	ctx context.Context
}

// Initialize a new per request public key store. Must be called before use
func NewPerRequestPublicKeyStore(endpoint string, ctx context.Context, opts ...PublicKeyStoreOpt) (*PerRequestPublicKeyStore, error) {
	s := &PerRequestPublicKeyStore{
		BasePublicKeyStore: BasePublicKeyStore{
			Endpoint: endpoint + "/api/pkeys",
			httpClient: http.DefaultClient,
		},
	}
	s.BasePublicKeyStore.keyFunc = s.keyFunc
	s.ctx = ctx

	for _, opt := range opts {
		if err := opt(&s.BasePublicKeyStore); err != nil {
			return nil, err
		}
	}

  if err := s.refreshPublicKeys(s.ctx); err != nil {
		return nil, err
	}
	return s, nil
}

// Checks and refreshes the keystore.
func (s *PerRequestPublicKeyStore) keyFunc(token *jwt.Token) (interface{}, error) {
	if time.Now().After(s.nextUpdate) {
		// Must unlock the keySetMutex because it is locked when coming into this function
		s.keySetMutex.RUnlock()
		s.refreshPublicKeys(s.ctx)
		s.keySetMutex.RLock()
	}
	return s.keySet, nil
}
