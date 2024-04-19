package stoke

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type WebCachePublicKeyStore struct {
	Endpoint string
	BasePublicKeyStore
}

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
