package stoke

import (
	"context"
	"net/http"
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
func NewWebCachePublicKeyStore(endpoint string, ctx context.Context, opts ...PublicKeyStoreOpt) (*WebCachePublicKeyStore, error) {
	s := &WebCachePublicKeyStore{
		BasePublicKeyStore: BasePublicKeyStore{
			Endpoint: endpoint,
			httpClient: http.DefaultClient,
		},
	}
	s.keyFunc = func(token *jwt.Token) (interface{}, error) {
		return s.keySet, nil
	}

	for _, opt := range opts {
		if err := opt(&s.BasePublicKeyStore); err != nil {
			return nil, err
		}
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
