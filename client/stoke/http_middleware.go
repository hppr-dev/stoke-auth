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
	ctx := req.Context()
	_, span := getTracer().Start(ctx, "AuthHandler.ServeHTTP")
	defer span.End()

	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Token ") {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(`{"message" : "Token is required"}`))
		return
	}

	trimToken := strings.TrimPrefix(token, "Token ")

	jwtToken, err := w.store.ParseClaims(ctx, trimToken, w.reqClaims, w.parserOpts...)
	if err != nil || 
			addTokenToSpan(jwtToken, span) || // This is a shortcut to always add the token to the span
			!jwtToken.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{"message" : "Unauthorized"}`))
		return
	}

	reqContext := context.WithValue(ctx, "token", trimToken)
	reqContext = context.WithValue(reqContext, "jwt.Token", jwtToken)

	w.inner(res, req.WithContext(reqContext))
}

func Token(ctx context.Context) *jwt.Token {
	return ctx.Value("jwt.Token").(*jwt.Token)
}

func TokenString(ctx context.Context) string {
	return ctx.Value("token").(string)
}
