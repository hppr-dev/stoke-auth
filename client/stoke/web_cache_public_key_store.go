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

// Initialize the WebCachePublicKeyStore. Must be called before use.
// Starts the management go routine.
func (s *WebCachePublicKeyStore) Init(ctx context.Context) error {
  s.BasePublicKeyStore.keyFunc = func(token *jwt.Token) (interface{}, error) {
		return s.keySet, nil
	}
	s.BasePublicKeyStore.Endpoint = s.Endpoint
  if err := s.refreshPublicKeys(ctx); err != nil {
		return err
	}
	go s.goManage(ctx)
	return nil
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
