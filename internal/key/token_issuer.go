package key

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	CustomClaims map[string]interface{} `json:"stk"`
}

type TokenIssuer interface {
	Init()
	IssueToken(claim Claims) (string, error)
	PublicString() string
}

type AsymetricTokenIssuer[P PrivateKey, K PublicKey]  struct {
	KeyPair[P, K]
}

func (a *AsymetricTokenIssuer[P, K]) IssueToken(claims Claims) (string, error) {
	priv, _ := a.Keys()
	return jwt.NewWithClaims(a.SigningMethod(), claims).SignedString(priv)
}

func (a *AsymetricTokenIssuer[P, K]) Init() {
	a.Generate()
}

func (a *AsymetricTokenIssuer[P, K]) PublicString() string {
	return a.KeyPair.PublicString()
}
