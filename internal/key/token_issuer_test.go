package key_test

import (
	"encoding/base64"
	"encoding/json"
	"crypto/ed25519"
	"stoke/internal/key"
	"stoke/internal/testutil"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

func TestAsymetricIssueTokenHappy(t *testing.T) {
	later := time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)
	earlier := time.Date(1000, time.January, 10, 10, 10, 10, 10, time.UTC)
	ctx := testutil.NewMockContext()

	claims := &stoke.Claims{
		StokeClaims: map[string]string {
			"hello" : "world",
			"foo": "bar",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Me",
			Subject:   "Myself",
			Audience:  []string{"eye"},
			ExpiresAt: jwt.NewNumericDate(later),
			NotBefore: jwt.NewNumericDate(earlier),
			IssuedAt:  jwt.NewNumericDate(earlier),
			ID:        "k1",
		},
	}

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	token, refresh, err := issuer.IssueToken(claims, ctx)
	if err != nil {
		t.Logf("An error occurred while generating token: %v", err)
		t.Fail()
	}

	expectToken(token, refresh, t)
}

func TestAsymetricIssueTokenWithTokenLimit(t *testing.T) {
	later := time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)
	earlier := time.Date(1000, time.January, 10, 10, 10, 10, 10, time.UTC)
	ctx := testutil.NewMockContext()

	claims := &stoke.Claims{
		StokeClaims: map[string]string {
			"hello" : "world",
			"foo": "bar",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Me",
			Subject:   "Myself",
			Audience:  []string{"eye"},
			ExpiresAt: jwt.NewNumericDate(later),
			NotBefore: jwt.NewNumericDate(earlier),
			IssuedAt:  jwt.NewNumericDate(earlier),
		},
	}

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		TokenRefreshLimit : 3,
		KeyCache: &MockKeyCache{},
	}

	token, refresh, err := issuer.IssueToken(claims, ctx)
	if err != nil {
		t.Logf("An error occurred while generating token: %v", err)
		t.Fail()
	}

	expectToken(token, refresh, t)
}

func TestAsymetricIssueTokenWithTokenLimitWithCustomKey(t *testing.T) {
	later := time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)
	earlier := time.Date(1000, time.January, 10, 10, 10, 10, 10, time.UTC)
	ctx := testutil.NewMockContext()

	claims := &stoke.Claims{
		StokeClaims: map[string]string {
			"ref": "k3",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Me",
			Subject:   "Myself",
			Audience:  []string{"eye"},
			ExpiresAt: jwt.NewNumericDate(later),
			NotBefore: jwt.NewNumericDate(earlier),
			IssuedAt:  jwt.NewNumericDate(earlier),
			ID:        "someID",
		},
	}

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		TokenRefreshLimit : 3,
		TokenRefreshCountKey: "ref",
		KeyCache: &MockKeyCache{},
	}

	token, refresh, err := issuer.IssueToken(claims, ctx)
	if err != nil {
		t.Logf("An error occurred while generating token: %v", err)
		t.Fail()
	}

	expectToken(token, refresh, t)
}

func TestAsymetricIssueTokenAtRefreshLimit(t *testing.T) {
	later := time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)
	earlier := time.Date(1000, time.January, 10, 10, 10, 10, 10, time.UTC)
	ctx := testutil.NewMockContext()

	claims := &stoke.Claims{
		StokeClaims: map[string]string {
			"ref": "k0",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Me",
			Subject:   "Myself",
			Audience:  []string{"eye"},
			ExpiresAt: jwt.NewNumericDate(later),
			NotBefore: jwt.NewNumericDate(earlier),
			IssuedAt:  jwt.NewNumericDate(earlier),
			ID:        "someID",
		},
	}

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		TokenRefreshLimit : 3,
		TokenRefreshCountKey: "ref",
		KeyCache: &MockKeyCache{},
	}

	token, refresh, err := issuer.IssueToken(claims, ctx)
	if err == nil {
		t.Logf("Was able to refresh token with reached token limit:\nt: %s\nr:%s", token, refresh)
		t.Fail()
	}
}

func TestAsymetricRefreshHappy(t *testing.T) {
	ctx := testutil.NewMockContext()

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"
	refreshStr := "02CwLlpkhkT2UHqanet2vztwchNC_7WhiwwJ2pK4sdd_FpBs4lJvWTTaCKvJARs3q0SkALJWLPYzfkW-pXkvDg=="

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &stoke.Claims{}, func(*jwt.Token) (interface{}, error) { return issuer.CurrentKey().PublicKey(), nil })
	if err != nil {
		t.Logf("An error occurent while parsing static token: %v", err)
		t.Fail()
	}

	token, refresh, err := issuer.RefreshToken(jwtToken, refreshStr, time.Hour, ctx)
	if err != nil {
		t.Logf("An error occurred while generating token: %v", err)
		t.Fail()
	}

	expectToken(token, refresh, t)
}

func TestAsymetricRefreshTokenBadBase64Encoding(t *testing.T) {
	ctx := testutil.NewMockContext()

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"
	refreshStr := "i***@@@##$%%%^"

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &stoke.Claims{}, func(*jwt.Token) (interface{}, error) { return issuer.CurrentKey().PublicKey(), nil })
	if err != nil {
		t.Logf("An error occurent while parsing static token: %v", err)
		t.Fail()
	}

	token, refresh, err := issuer.RefreshToken(jwtToken, refreshStr, time.Hour, ctx)
	if err == nil {
		t.Logf("Was able to refresh with bad base64 refresh token:\nt: %s\nr: %s", token, refresh)
		t.Fail()
	}
}

func TestAsymetricRefreshTokenInvalidRefreshToken(t *testing.T) {
	ctx := testutil.NewMockContext()

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"
	refreshStr := "abcdefGHIJKabc=="

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &stoke.Claims{}, func(*jwt.Token) (interface{}, error) { return issuer.CurrentKey().PublicKey(), nil })
	if err != nil {
		t.Logf("An error occurent while parsing static token: %v", err)
		t.Fail()
	}

	token, refresh, err := issuer.RefreshToken(jwtToken, refreshStr, time.Hour, ctx)
	if err == nil {
		t.Logf("Was able to refresh with unverifiable refresh token:\nt: %s\nr: %s", token, refresh)
		t.Fail()
	}
}

func TestAsymetricRefreshTokenInvalidClaimsType(t *testing.T) {
	ctx := testutil.NewMockContext()

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"
	refreshStr := "02CwLlpkhkT2UHqanet2vztwchNC/7WhiwwJ2pK4sdd/FpBs4lJvWTTaCKvJARs3q0SkALJWLPYzfkW+pXkvDg=="

	mapClaims := make(jwt.MapClaims)
	jwtToken, err := jwt.ParseWithClaims(tokenStr, mapClaims, func(*jwt.Token) (interface{}, error) { return issuer.CurrentKey().PublicKey(), nil })
	if err != nil {
		t.Logf("An error occurent while parsing static token: %v", err)
		t.Fail()
	}

	token, refresh, err := issuer.RefreshToken(jwtToken, refreshStr, time.Hour, ctx)
	if err == nil {
		t.Logf("Was able to refresh with bad claims type:\nt: %s\nr: %s", token, refresh)
		t.Fail()
	}
}

func TestAsymetricWithContextFromContext(t *testing.T) {
	ctx := testutil.NewMockContext()

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	newCtx := issuer.WithContext(ctx)

	ctxIssuer := key.IssuerFromCtx(newCtx)

	if ctxIssuer.(*key.AsymetricTokenIssuer[ed25519.PrivateKey]) != &issuer {
		t.Log("Issuer from context is not the same as the one inserted")
		t.Fail()
	}

}

func expectToken(token, refresh string, t *testing.T) {
	splitToken := strings.Split(token, ".")
	if len(splitToken) != 3 {
		t.Fatalf("token should have 3 parts, got %d", len(splitToken))
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(splitToken[0])
	if err != nil {
		t.Fatalf("decode header: %v", err)
	}
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		t.Fatalf("parse header: %v", err)
	}
	if header["alg"] != "EdDSA" || header["typ"] != "JWT" {
		t.Fatalf("header alg/typ: want EdDSA/JWT, got %v", header)
	}
	if header["kid"] != "p-0" {
		t.Fatalf("header kid: want p-0, got %v", header["kid"])
	}

	bodyBytes, err := base64.RawURLEncoding.DecodeString(splitToken[1])
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		t.Fatalf("parse payload: %v", err)
	}
	if payload["kid"] != "p-0" {
		t.Fatalf("payload kid: want p-0, got %v", payload["kid"])
	}

	// Refresh token must be non-empty valid base64 (EdDSA refresh is 64 bytes)
	if refresh == "" {
		t.Fatal("refresh token is empty")
	}
	refreshBytes, err := base64.URLEncoding.DecodeString(refresh)
	if err != nil {
		t.Fatalf("refresh token not valid base64: %v", err)
	}
	if len(refreshBytes) != 64 {
		t.Fatalf("refresh token length: want 64 (EdDSA), got %d", len(refreshBytes))
	}
}

