package web

import (
	"log"
	"net/http"
	"stoke/internal/ctx"
	"stoke/internal/key"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginApiHandler struct {
	Context *ctx.Context
}

func (l LoginApiHandler) ServeHTTP(res http.ResponseWriter, _ *http.Request) {
	claims := key.Claims{
		CustomClaims : map[string]interface{} {
			"hlo" : "world",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "stoke",
			Subject: "auth",
			Audience: []string{ "yomoma" },
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID: "stoke-1",
		},
	}
	token, err := l.Context.Issuer.IssueToken(claims)
	if err != nil {
		log.Printf("Uhoh: %v", err)
	}
	res.Write([]byte(token))
}

type PkeyApiHandler struct {
	Context *ctx.Context
}

func (p PkeyApiHandler) ServeHTTP(res http.ResponseWriter, _ *http.Request) {
	b, err := p.Context.Issuer.PublicKeys()
	if err != nil {
		log.Printf("Could not get public keys: %v", err)
	}
	res.Write(b)
}
