package stoke

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type WebPublicKeyStore struct {
	Endpoint   string
	nextUpdate time.Time
	keySet jwt.VerificationKeySet
}

func (s *WebPublicKeyStore) Init() error {
  if err := s.refreshPublicKeys(); err != nil {
		return err
	}
	go s.goManage()
	return  nil
}

func (s *WebPublicKeyStore) goManage(){
	for {
		select {
		case <-time.After(time.Now().Sub(s.nextUpdate)):
			s.refreshPublicKeys()
		}
	}
}

func (s *WebPublicKeyStore) ParseClaims(token string, reqClaims *Claims, parserOpts ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, reqClaims, s.keyFunc, parserOpts...)
}

func (s *WebPublicKeyStore) keyFunc(token *jwt.Token) (interface{}, error) {
	// TODO Test performance impact of checking key validity here instead of go routine
	return s.keySet, nil
}

func (s *WebPublicKeyStore) refreshPublicKeys() error {
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

	s.keySet = jwt.VerificationKeySet {
		Keys : pkeys,
	}
	s.nextUpdate = jwks.Expires
	return nil
}
