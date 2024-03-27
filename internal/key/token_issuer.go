package key

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	CustomClaims map[string]interface{} `json:"stk"`
}

type TokenIssuer interface {
	IssueToken(claim Claims) (string, error)
	PublicKeys() ([]byte, error)
}

type AsymetricTokenIssuer[P PrivateKey]  struct {
	KeyCache[P]
}

func (a *AsymetricTokenIssuer[P]) IssueToken(claims Claims) (string, error) {
	curr := a.CurrentKey()
	return jwt.NewWithClaims(curr.SigningMethod(), claims).SignedString(curr.Key())
}

func (a *AsymetricTokenIssuer[P]) PublicKeys() ([]byte, error) {
	return a.JSON()
}
