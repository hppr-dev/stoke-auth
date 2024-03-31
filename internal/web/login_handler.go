package web

import (
	"errors"
	"fmt"
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

func (l LoginApiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		MethodNotAllowed.Write(res)
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
		BadRequest.Write(res)
		return
	}
	claims, err := l.Context.UserProvider.GetUserClaims(username, password)
	if err != nil {
		Unauthorized.Write(res)
		return
	}
	
	claimMap := make(map[string]string)

	for _, claim := range claims {
		claimMap[claim.ShortName] = claim.Value
	}

	token, err := l.Context.Issuer.IssueToken(key.Claims{
		StokeClaims : claimMap,
		RegisteredClaims: l.createRegisteredClaims(),
	})
	if err != nil {
		InternalServerError.Write(res)
		return
	}

	res.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token)))
}

func (l LoginApiHandler) createRegisteredClaims() jwt.RegisteredClaims {
	now := time.Now()
	return jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 30)),
		Issuer:    l.Context.Config.Tokens.Issuer,
		Subject:   l.Context.Config.Tokens.Subject,
		Audience:  l.Context.Config.Tokens.Audience,
	}
}
