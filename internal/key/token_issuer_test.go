package key_test

import (
	"crypto/ed25519"
	"stoke/internal/key"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

func TestAsymetricIssueTokenHappy(t *testing.T) {
	later := time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)
	earlier := time.Date(1000, time.January, 10, 10, 10, 10, 10, time.UTC)
	ctx := NewMockContext()

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

	// We need to handle any valid ordering of keys and order will change the signature
	expBody := map[string]string {
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ" : "02CwLlpkhkT2UHqanet2vztwchNC/7WhiwwJ2pK4sdd/FpBs4lJvWTTaCKvJARs3q0SkALJWLPYzfkW+pXkvDg==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJmb28iOiJiYXIiLCJoZWxsbyI6IndvcmxkIn0.yBC1uvZqMm9OO_elMqOudP8tFbpMeJ8Q8YxejoQ14Cay3pH0I_qHoc0r7bNVO3bH99aVrqVYTtdzZbeJ9wymBA" : "VY94KsCKyBcYtXx457F+MGmP93T16OZg6U3UJmP/l6//ruMaAUHfYfNdBzsDrXItg292ENhJXBedhy1CAN4ZCg==",
	}

	expectToken(token, refresh, expBody, t)
}

func TestAsymetricRefreshHappy(t *testing.T) {
	ctx := NewMockContext()

	issuer := key.AsymetricTokenIssuer[ed25519.PrivateKey]{
		KeyCache: &MockKeyCache{},
	}

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"
	refreshStr := "02CwLlpkhkT2UHqanet2vztwchNC/7WhiwwJ2pK4sdd/FpBs4lJvWTTaCKvJARs3q0SkALJWLPYzfkW+pXkvDg=="

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

	// We need to handle any valid ordering of keys and order will change the signature
	expBodyMap := map[string]string {
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODQwMTgxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.tL4uWbwk1sMJB6WxszWZM26E1CEtIOSSPemQLJqAyxlsno6i2saKaCcdlC1Iy4WuAq2NiKB8sZUMmgeyRdLWCw": "7aQHgfk8jADRErR+EuaAtSxx19kAZEZxAdY8zRL6OKGlpgyJK30Udw0EW75OCgzbAKbu8WUPyTLTMZ6tnABMBA==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODQwMTgxMCwianRpIjoiazEiLCJmb28iOiJiYXIiLCJoZWxsbyI6IndvcmxkIn0.-FU1dKzTjzMOEXDznuOs2U5KGVb_9jKJFOAnHGFZZerewVsbkYdFa7GOTs9qorNAAaX--A9K_GyiXEKDrSZaCg": "4o5X9DQe/cqILBeu983enrdlJeXny0d/l78Tp5VXoGxYxIEWIOtTKWaTBH7voObNS5RtSFyPlb8GME6tXS2ZDA==",
	}

	expectToken(token, refresh, expBodyMap, t)
}

func TestAsymetricWithContextFromContext(t *testing.T) {
	ctx := NewMockContext()

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

func expectToken(token, refresh string, expBodyMap map[string]string, t *testing.T) {
	expHeader := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9"

	splitToken := strings.Split(token, ".")
	if splitToken[0] != expHeader {
		t.Logf("Header did not match expected value:\nE:%s\nA:%s", expHeader, splitToken[0])
		t.Fail()
	}

	bodySig := splitToken[1] + "." + splitToken[2]
	expRefresh, ok := expBodyMap[bodySig]
	if !ok {
		t.Logf("Token Body.Signature was not recognized:\n\n\"%s\" : \"%s\",\n\n", bodySig, refresh)
		t.FailNow()
	}

	if refresh != expRefresh {
		t.Logf("Refresh token did not match expected:\nT:%s\n\nE:%s\nA:%s", token, expRefresh, refresh)
		t.Fail()
	}
}
