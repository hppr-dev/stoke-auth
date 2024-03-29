package stoke

import (
	"net/http"
	"strings"
)

func WithClaims(handler http.Handler, store PublicKeyStore, claims *ClaimsValidator) http.Handler {
	return authWrapper{
		inner: handler.ServeHTTP,
		store: store,
		reqClaims : claims,
	}
}

func WithClaimsFunc(handler http.HandlerFunc, store PublicKeyStore, claims *ClaimsValidator) http.Handler {
	return authWrapper{
		inner:     handler,
		store:     store,
		reqClaims: claims,
	}
}

type authWrapper struct {
	store PublicKeyStore
	inner http.HandlerFunc
	reqClaims *ClaimsValidator
}

func (w authWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Token ") || !w.store.ValidateClaims(strings.TrimPrefix(token, "Token "), w.reqClaims) {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("{'message' : 'Unauthorized'}"))
		return
	}
	w.inner(res, req)
}

