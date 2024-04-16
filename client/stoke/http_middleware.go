package stoke

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(handler http.Handler, store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler.ServeHTTP,
		TokenHandler: NewTokenHandler(store, claims, parserOpts...),
	}
}

func AuthFunc(handler http.HandlerFunc, store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler,
		TokenHandler: NewTokenHandler(store, claims, parserOpts...),
	}
}

type authWrapper struct {
	inner http.HandlerFunc
	*TokenHandler
}

func (w authWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx, span := getTracer().Start(ctx, "AuthHandler.ServeHTTP")
	defer span.End()

	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(`{"message" : "Token is required"}`))
		return
	}

	trimToken := strings.TrimPrefix(token, "Bearer ")

	ctx, err := w.InjectToken(trimToken, ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{"message" : "Unauthorized"}`))
		return
	}

	w.inner(res, req.WithContext(ctx))
}
