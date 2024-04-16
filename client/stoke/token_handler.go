package stoke

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	InvalidTokenError error = errors.New("Invalid token")
)

func NewTokenHandler(store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) *TokenHandler {
	return &TokenHandler{
		store: store,
		reqClaims: claims,
		parserOpts: parserOpts,
	}
}

type TokenHandler struct {
	store PublicKeyStore
	reqClaims *Claims
	parserOpts []jwt.ParserOption
}

// Parses and injects the parsed token into the given context
// Returns non-nil error on an invalid token
func (w TokenHandler) InjectToken(token string, ctx context.Context) (context.Context, error) {
	_, span := getTracer().Start(ctx, "TokenHandler.InjectToken")
	defer span.End()

	if token == "" {
		return ctx, InvalidTokenError
	}

	jwtToken, err := w.store.ParseClaims(ctx, token, w.reqClaims, w.parserOpts...)
	if err != nil || 
			AddTokenToSpan(jwtToken, span) || // This is a shortcut to always add the token to the span
			!jwtToken.Valid {
		return ctx, InvalidTokenError
	}

	return context.WithValue(ctx, "jwt.Token", jwtToken), nil

}

func Token(ctx context.Context) *jwt.Token {
	return ctx.Value("jwt.Token").(*jwt.Token)
}
