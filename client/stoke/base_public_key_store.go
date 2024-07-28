package stoke

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// BasePublicKeyStore implements contains common public key store functionality
type BasePublicKeyStore struct {
	Endpoint   string
	httpClient *http.Client
	nextUpdate time.Time
	keySet jwt.VerificationKeySet
	keyFunc jwt.Keyfunc
	keySetMutex sync.RWMutex
}

// Parses Claims according to the configured required claims and parser options.
// keyFuncs must RUnlock before refreshPublicKeys
func (s *BasePublicKeyStore) ParseClaims(ctx context.Context, token string, reqClaims *Claims, parserOpts ...jwt.ParserOption) (*jwt.Token, error) {
	_, span := getTracer().Start(ctx, "ClientKeyStore.ParseClaims")
	defer span.End()

	s.keySetMutex.RLock()
	defer s.keySetMutex.RUnlock()

	return jwt.ParseWithClaims(token, reqClaims.New(), s.keyFunc, parserOpts...)
}


// Sets the tls config for the key store
func (s *BasePublicKeyStore) SetTLSConfig(tlsConfig *tls.Config) {
	s.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}


// refreshPublicKeys retreives public keys from the configured endpoint and saves them to the store
func (s *BasePublicKeyStore) refreshPublicKeys(ctx context.Context) error {
	_, span := getTracer().Start(ctx, "ClientKeyStore.refreshPublicKeys")
	defer span.End()

	resp, err := s.httpClient.Get(s.Endpoint)
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

	s.keySetMutex.Lock()
	defer s.keySetMutex.Unlock()

	s.keySet = jwt.VerificationKeySet {
		Keys : pkeys,
	}
	s.nextUpdate = jwks.Expires

	return nil
}
