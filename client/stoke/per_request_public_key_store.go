package stoke

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PerRequestPublicKeyStore struct {
	Endpoint   string
	BasePublicKeyStore
	ctx context.Context
}

func (s *PerRequestPublicKeyStore) Init(ctx context.Context) error {
	s.BasePublicKeyStore.Endpoint = s.Endpoint
	s.BasePublicKeyStore.keyFunc = s.keyFunc
	s.ctx = ctx
  if err := s.refreshPublicKeys(s.ctx); err != nil {
		return err
	}
	return nil
}

func (s *PerRequestPublicKeyStore) keyFunc(token *jwt.Token) (interface{}, error) {
	if time.Now().After(s.nextUpdate) {
		s.refreshPublicKeys(s.ctx)
	}
	return s.keySet, nil
}
