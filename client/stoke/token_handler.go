package stoke

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	InvalidTokenError error = errors.New("Invalid token")
)

// Creates a new TokenHandler that verifies claims against a public keystore
func NewTokenHandler(store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) *TokenHandler {
	return &TokenHandler{
		store: store,
		reqClaims: claims,
		parserOpts: parserOpts,
	}
}

// TokenHandler is resposible for verifying token claims against a public key store 
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

	jwtToken, err := w.store.ParseClaims(ctx, token, w.reqClaims.New(), w.parserOpts...)
	if err != nil || 
			AddTokenToSpan(jwtToken, span) || // This is a shortcut to always add the token to the span
			!jwtToken.Valid {
		return ctx, InvalidTokenError
	}

	return context.WithValue(ctx, "jwt.Token", jwtToken), nil

}

// Gets the current JWT from the context
func Token(ctx context.Context) *jwt.Token {
	return ctx.Value("jwt.Token").(*jwt.Token)
}
