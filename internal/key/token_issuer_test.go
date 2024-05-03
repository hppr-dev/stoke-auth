package key_test

import (
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

	// We need to handle any valid ordering of keys and order will change the signature
	expBody := map[string]string {
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIiwia2lkIjoiMCJ9.P74Wzd0mBlPXz6kpYTna6ud1F1_ngu5BsTHQq1FpgLai2jiOb2KKMo4G9BEMy8uNkjqjQzJzR7gGQeiz_lBRDw" : "fGha64MbiPNyAUgh19F2-GUXBckujT4O6CBMpHSQ9oAvo2fqHLoxuzSoR0TIv7Q8o_nc0iqqySAMGwj2PBuDCA==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJmb28iOiJiYXIiLCJraWQiOiIwIiwiaGVsbG8iOiJ3b3JsZCJ9.rMo-hVR1fXsV1JiWzNHPCrkDKXhT6-R32Rs3pBNkNRees2w958fV_W3Y0XXR3lwPBUyf1JsFrqMFdLxXzI6cDA" : "v5O55dwC7kr_E8i24pRyJcX7O37g0_3LOeEsGrmWchQ7yy-zHwQSKWrNy-5zx5mtprDQqNYY5CM1BBc2UhdUBw==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJraWQiOiIwIiwiaGVsbG8iOiJ3b3JsZCIsImZvbyI6ImJhciJ9.GKAuf1wlBv0uL4uYUSWoAQwnajPNK_n6lADcegmiTiJbVcpz2I8RYxUE7EZna7N-d4NjSY2zKVuainwYDy4NBw" : "I72ZvdHoQmEdnzCdfTA9SE9RfWmqyG1hl4fQ8-MV4qId_Lx8nDiFumHHzuYTDj8aal-Z7rk_Ji_XMqgdXC4rDQ==",
	}

	expectToken(token, refresh, expBody, t)
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

	// We need to handle any valid ordering of keys and order will change the signature
	expBody := map[string]string {
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazMiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIiwia2lkIjoiMCJ9.J-aaXAWsCGTvL6TTwtRp8g3TOWjZoLbIBqihEEaW3JnBZSetg9vzWbXOhBI3uXygCFUd_zNS444CD2OUCoMIAw" : "kSFtOjwE9n40J6BthpUkmEVmkrDGOU02gQJAje9UA6QyxAEQ2PUKK2ch58qAEYuu8ry8D1WaYwUID9OSkXFHCg==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazMiLCJraWQiOiIwIiwiaGVsbG8iOiJ3b3JsZCIsImZvbyI6ImJhciJ9.PcCoNF3l24HBTZxTisSawYtubygTw2nzjvZ5km0w4fdhhBltBqSanCbPh_g_mKDYVvq063Im6Qk2w9n_xEYKBA" : "baUX214V35xO2D9m1LzgTLz7n_4bT8d1miZBHHMCjFVGkK2v_w06tZJQOsnXuPPWqGnlw0V6M545Zd006SIDAA==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazMiLCJmb28iOiJiYXIiLCJraWQiOiIwIiwiaGVsbG8iOiJ3b3JsZCJ9.g-25mWypad_Lj_TNU4JHMWlhUXHzISO_D5RWj00tM25XAQFUaJQT7eXbhQXrqCzKF0bJy31IT0Y1ZkqMPyd_Bw" : "qzQWJQHo0hO4GlesFsiszR4VBa8Ze_0_NIeJglPaEXq0lXpbDsJEEyEn_0XonsIHnOZBANEYeVOF6hyzuAY4AQ==",
	}

	expectToken(token, refresh, expBody, t)
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

	// We need to handle any order of keys
	expBody := map[string]string {
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoic29tZUlEIiwicmVmIjoiazIiLCJraWQiOiIwIn0.KlH-j6BcRltIBcQH3FkaoBFS0Dy5RgSpgp0clVbArSaegyiJinEyntrfVW3fOgm94SnyCstuwCJ6LKOg3GQKDA" : "G9z2rMfOqUTyj20faiqf5J5S_2eHmm2iO1M1IU8-4ageumx2Me6ML-sAWxR_fHH3J6-_AbSvU-iZuWbPeVp5Cw==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoic29tZUlEIiwia2lkIjoiMCIsInJlZiI6ImsyIn0.2tdrXw3GsjYOEQTetz3lWs5k8N-qZb5TcnWRj0B_TaxwUijbxH-L46lRT2nBHC5VSo_Ld4RQkVlD2qlTkm9fCw" : "3zi6qG1TQSbbnpZE8whQSpcMcL6FjH6g0wu1M5II309M-0ge6qyV_2l5V3mlfi4BG10uAvmnID0SsSDXwvLMBw==",
	}

	expectToken(token, refresh, expBody, t)
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

	// We need to handle any valid ordering of keys and order will change the signature
	expBodyMap := map[string]string {
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODQwMTgxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIiwia2lkIjoiMCJ9.UARbF9h-E_95vxMn1k4Gk4AmeJIak6hLCfb-VdJ534MRtRcu1pGoXO6LMpJdLDIdKwJtTt2fBaj9bB5HjNzQBw" : "4995wpo_tQer0QRLHiGAwRXQ6WgXgn_DeRom8Q8102T3m8be0I7ZiIUGP2iR1UoI6Ms1jazAQJx8qJWAst9zAw==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODQwMTgxMCwianRpIjoiazEiLCJraWQiOiIwIiwiaGVsbG8iOiJ3b3JsZCIsImZvbyI6ImJhciJ9.TdsL5RtWxUpkkvDiP5ZDkLyxvYhERrKnmTqLt55kIQQ6azrQpmTFnp8VgRo7WIyieMGPh2Rvlk9Y0YjfHO9lAQ" : "lHc_M5CT12mFmM9p49dUHKmRBy3j6AE7nZWeylt3IcZJULok2l6FUxZ3sMj8d0okJab6P1EUXntesAxtXVTLDg==",
		"eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODQwMTgxMCwianRpIjoiazEiLCJmb28iOiJiYXIiLCJraWQiOiIwIiwiaGVsbG8iOiJ3b3JsZCJ9.07x6XWC5vtk-Jc6MuxL8N7Hc5xZA2DEtlIqNmbfmzXE22s6RYRj372kj8H6OGQen6_a2bzjYWP1uw1UHKWygBw" : "062974mDze4qFldjRB2ySXAe6_VyorcTithEpWzD45obgqWD2jiSCR9D0ppBoujnE5VDMiOacCFfYTUDYZ4uAA==",
	}

	expectToken(token, refresh, expBodyMap, t)
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

func expectToken(token, refresh string, expBodyMap map[string]string, t *testing.T) {
	expHeader := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9"

	splitToken := strings.Split(token, ".")
	if splitToken[0] != expHeader {
		t.Logf("Header did not match expected value:\nE:%s\nA:%s", expHeader, splitToken[0])
		t.FailNow()
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

