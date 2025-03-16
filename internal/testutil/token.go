package testutil

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

type TokenOption func(mockTokenHandler) mockTokenHandler

// Adds a token to the handler
func WithToken(t *testing.T, opts ...TokenOption) ContextOption {
	return func(ctx context.Context) context.Context {
		mth := mockTokenHandler{}
		mth.setDefaults()
		for _, opt := range opts {
			mth = opt(mth)
		}
		mth.reinitialize()

		ctx, err := mth.handler.InjectToken(mth.rawToken, ctx)
		if err != nil {
			t.Fatalf("Could not inject token!")
		}
		return ctx
	}
}

type mockTokenHandler struct {
	rawToken       string
	defaultToken   *jwt.Token
	keystore       *stoke.TestPublicKeyStore
	requiredClaims *stoke.Claims
	handler        *stoke.TokenHandler
}

func WithTokenClaim(key, value string) TokenOption {
	return func (m mockTokenHandler) mockTokenHandler {
		claims, _ := m.defaultToken.Claims.(jwt.MapClaims)
		claims[key] = value
		m.defaultToken.Claims = claims
		return m
	}
}

func WithRawToken(rawToken string) TokenOption {
	return func (m mockTokenHandler) mockTokenHandler {
		m.rawToken = rawToken
		return m
	}
}

func (m *mockTokenHandler) setDefaults() {
	m.rawToken = ""
	m.defaultToken = jwt.New(jwt.SigningMethodNone)
}

func (m *mockTokenHandler) reinitialize() {
	m.keystore = stoke.NewTestPublicKeyStore(m.defaultToken)
	m.keystore.SetAllow()
	m.requiredClaims = stoke.RequireToken()
	m.handler = stoke.NewTokenHandler(m.keystore, m.requiredClaims)
}
