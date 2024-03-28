package web

import (
	"errors"
	"log"
	"net/http"
	"stoke/internal/ctx"
	"stoke/internal/key"
	"time"

	"github.com/go-faster/jx"
	"github.com/golang-jwt/jwt/v5"
)

type LoginApiHandler struct {
	Context *ctx.Context
}

type loginRequest struct {
	Username string
	Password string
}

func (l LoginApiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte("Method Not Allowed"))
		return
	}

	var username, password string
	decoder := jx.Decode(req.Body, 256)
	err := decoder.Obj(func (d *jx.Decoder, key string) error {
		var err error
		switch key {
		case "username":
			username, err = d.Str()
		case "password":
			password, err = d.Str()
		default:
			return errors.New("Bad Request")
		}
		return err
	})

	if err != nil || username == "" || password == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Bad Request: Missing parameters"))
		return
	}
	claims, err := l.Context.UserProvider.GetUserClaims(username, password)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Unauthorized"))
		return
	}
	
	claimMap := make(map[string]string)

	for _, claim := range claims {
		claimMap[claim.ShortName] = claim.Value
	}

	token, err := l.Context.Issuer.IssueToken(key.Claims{
		StokeClaims : claimMap,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "stk",
			Subject: "ath",
			Audience: []string{ "stkuser" },
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal Server Error"))
		return
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
