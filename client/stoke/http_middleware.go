package stoke

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(handler http.Handler, store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler.ServeHTTP,
		store:      store,
		reqClaims : claims,
		parserOpts: parserOpts,
	}
}

func AuthFunc(handler http.HandlerFunc, store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler,
		store:      store,
		reqClaims:  claims,
		parserOpts: parserOpts,
	}
}

type authWrapper struct {
	store PublicKeyStore
	inner http.HandlerFunc
	reqClaims *Claims
	parserOpts []jwt.ParserOption
}

func (w authWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Token ") {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(`{"message" : "Token is required"}`))
		return
	}

	trimToken := strings.TrimPrefix(token, "Token ")

	jwtToken, err := w.store.ParseClaims(trimToken, w.reqClaims, w.parserOpts...)
	if err != nil  || !jwtToken.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{"message" : "Unauthorized"}`))
		return
	}

	reqContext := context.WithValue(req.Context(), "token", trimToken)
	reqContext = context.WithValue(reqContext, "jwt.Token", jwtToken)

	w.inner(res, req.WithContext(reqContext))
}

