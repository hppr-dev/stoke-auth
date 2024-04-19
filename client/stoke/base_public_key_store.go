package stoke

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type BasePublicKeyStore struct {
	Endpoint   string
	nextUpdate time.Time
	keySet jwt.VerificationKeySet
	keyFunc jwt.Keyfunc
	mutex sync.Mutex
}

func (s *BasePublicKeyStore) ParseClaims(ctx context.Context, token string, reqClaims *Claims, parserOpts ...jwt.ParserOption) (*jwt.Token, error) {
	_, span := getTracer().Start(ctx, "ClientKeyStore.ParseClaims")
	defer span.End()

	// TODO check for client data races
	return jwt.ParseWithClaims(token, reqClaims.New(), s.wrappedKeyFunc, parserOpts...)
}

func (s *BasePublicKeyStore) wrappedKeyFunc(token *jwt.Token) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.keyFunc(token)
}

func (s *BasePublicKeyStore) refreshPublicKeys(ctx context.Context) error {
	_, span := getTracer().Start(ctx, "ClientKeyStore.refreshPublicKeys")
	defer span.End()

	resp, err := http.Get(s.Endpoint)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jwks := &JWKSet{}
	if err := json.Unmarshal(bodyBytes, jwks); err != nil {
		return err
	}

	pkeys := make([]jwt.VerificationKey, len(jwks.Keys))
	for i, k := range jwks.Keys {
		pkeys[i], err = k.ToPublicKey()
		if err != nil {
			return err
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.keySet = jwt.VerificationKeySet {
		Keys : pkeys,
	}
	s.nextUpdate = jwks.Expires

	return nil
}
