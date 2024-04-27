package stoke

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// WebCachePublicKeyStore periodically pulls the public keys from the stoke server.
// This type of keystore keeps the set of public keys as up to date as possible
type WebCachePublicKeyStore struct {
	Endpoint string
	BasePublicKeyStore
}

// Initialize a WebCachePublicKeyStore
// Starts the management go routine.
func NewWebCachePublicKeyStore(endpoint string, ctx context.Context) (*WebCachePublicKeyStore, error) {
	s := &WebCachePublicKeyStore{
		BasePublicKeyStore: BasePublicKeyStore{
			Endpoint: endpoint,
		},
	}
	s.keyFunc = func(token *jwt.Token) (interface{}, error) {
		return s.keySet, nil
	}

  if err := s.refreshPublicKeys(ctx); err != nil {
		return nil, err
	}

	go s.goManage(ctx)
	return s, nil
}

// Manages the keystore from a go routine, pulling public keys at the time specified by exp
func (s *WebCachePublicKeyStore) goManage(ctx context.Context){
	for {
		select {
		case <-ctx.Done():
			break
		case <-time.After(time.Now().Sub(s.nextUpdate)):
			s.refreshPublicKeys(ctx)
		}
	}
}
