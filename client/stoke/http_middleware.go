package stoke

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func WithClaims(handler http.Handler, store PublicKeyStore, claims *ClaimsValidator, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler.ServeHTTP,
		store:      store,
		reqClaims : claims,
		parserOpts: parserOpts,
	}
}

func WithClaimsFunc(handler http.HandlerFunc, store PublicKeyStore, claims *ClaimsValidator, parserOpts ...jwt.ParserOption) http.Handler {
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
	reqClaims *ClaimsValidator
	parserOpts []jwt.ParserOption
}

func (w authWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Token ") || !w.store.ValidateClaims(strings.TrimPrefix(token, "Token "), w.reqClaims, w.parserOpts...) {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("{'message' : 'Unauthorized'}"))
		return
	}
	w.inner(res, req)
}

