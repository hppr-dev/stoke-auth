package stoke

import (
	"net/http"
	"time"

	"github.com/go-faster/jx"
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

func (s *WebPublicKeyStore) ValidateClaims(token string, reqClaims *ClaimsValidator) bool {
	// TODO Make audience/issuer configurable
	jwtToken, err := jwt.ParseWithClaims(token, reqClaims, s.keyFunc, jwt.WithIssuer("stk"))
	if err != nil {
		return false
	}
	return jwtToken.Valid
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
	decoder := jx.Decode(resp.Body, 256)

	var pkeys []jwt.VerificationKey
	var nextUpdate time.Time

	err = decoder.Arr(
		func (d *jx.Decoder) error {
			return d.Obj(
				func (d *jx.Decoder, key string) error {
					var objErr error
					var newKey []uint8
					var t time.Time
					switch key {
					case "text":
						newKey, objErr = d.Base64()
						if objErr != nil {
							pkeys = append(pkeys, newKey)
						}
					case "renews", "expires":
						t, objErr = parseTime(d)
						if t.Before(nextUpdate) {
							nextUpdate = t
						}
					}
					return objErr
				},
			)
		},
	)
	if err != nil {
		return err
	}
	s.keySet = jwt.VerificationKeySet {
		Keys : pkeys,
	}
	s.nextUpdate = nextUpdate
	return nil
}

func parseTime(d *jx.Decoder) (time.Time, error) {
	i, err := d.Int()
	return time.Unix(int64(i), 0), err
}
