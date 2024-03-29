package key

import (
	"stoke/client/stoke"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	StokeClaims map[string]string `json:"stk"`
}

type TokenIssuer interface {
	IssueToken(claim Claims) (string, error)
	PublicKeys() ([]byte, error)
	stoke.PublicKeyStore
}

type AsymetricTokenIssuer[P PrivateKey]  struct {
	KeyCache[P]
}

func (a *AsymetricTokenIssuer[P]) IssueToken(claims Claims) (string, error) {
	curr := a.CurrentKey()
	return jwt.NewWithClaims(curr.SigningMethod(), claims).SignedString(curr.Key())
}
