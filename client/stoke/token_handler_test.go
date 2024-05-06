package stoke_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

func TestInjectValidToken(t *testing.T) {
	now := time.Now()
	defToken := jwt.NewWithClaims(jwt.SigningMethodNone, &stoke.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1234567890",
			IssuedAt: jwt.NewNumericDate(now),
		},
		StokeClaims: map[string]string{
			"name" : "John Doe",
		},
	})
	store := stoke.NewTestPublicKeyStore(defToken)
	handler := stoke.NewTokenHandler(store, stoke.RequireToken())

	ctx, err := handler.InjectToken("", context.Background())
	if err != nil {
		t.Logf("An error occured injecting token: %v", err)
		t.Fail()
	}

	token := stoke.Token(ctx)
	if claims, ok := token.Claims.(*stoke.Claims); !ok || claims.Subject != "1234567890" || claims.IssuedAt.Unix() != now.Unix() || claims.StokeClaims["name"] != "John Doe" {
		t.Logf("Claims did not match: %v", claims)
		t.Fail()
	}
}

func TestInjectInvalidToken(t *testing.T) {
	store := stoke.NewTestPublicKeyStore(nil)
	store.SetInvalid()
	handler := stoke.NewTokenHandler(store, stoke.RequireToken())

	if _, err := handler.InjectToken("invalidtoken", context.Background()); err == nil {
		t.Fatal("An error was not returned")
	}
}
