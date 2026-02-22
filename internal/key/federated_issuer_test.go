package key

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"stoke/internal/cluster"
	"hppr.dev/stoke"
)

// mockFederatedInner is a TokenIssuer that returns fixed JWKS from PublicKeys.
type mockFederatedInner struct {
	publicKeysBytes []byte
}

func (m *mockFederatedInner) IssueToken(*stoke.Claims, context.Context) (string, string, error) {
	return "", "", nil
}
func (m *mockFederatedInner) RefreshToken(*jwt.Token, string, time.Duration, context.Context) (string, string, error) {
	return "", "", nil
}
func (m *mockFederatedInner) PublicKeys(ctx context.Context) ([]byte, error) {
	return m.publicKeysBytes, nil
}
func (m *mockFederatedInner) WithContext(ctx context.Context) context.Context {
	return ctx
}
func (m *mockFederatedInner) ParseClaims(context.Context, string, *stoke.Claims, ...jwt.ParserOption) (*jwt.Token, error) {
	return nil, nil
}

func TestFederatedTokenIssuer_PublicKeys_ReturnsMergedJWKS(t *testing.T) {
	exp := time.Now().Add(time.Hour)
	localJWKS := stoke.JWKSet{
		Expires: exp,
		Keys: []*stoke.JWK{
			{KeyId: "p-0", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x0", Y: "y0"},
		},
	}
	localBytes, err := json.Marshal(localJWKS)
	if err != nil {
		t.Fatal(err)
	}

	peerJWKS := stoke.JWKSet{
		Expires: exp,
		Keys: []*stoke.JWK{
			{KeyId: "p-1", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x1", Y: "y1"},
		},
	}
	peerBytes, err := json.Marshal(peerJWKS)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pkeys" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(peerBytes)
	}))
	defer srv.Close()

	inner := &mockFederatedInner{publicKeysBytes: localBytes}
	discoverer := &cluster.StaticDiscoverer{URLs: []string{srv.URL}}
	federated := NewFederatedTokenIssuer(inner, discoverer, nil, "", 0)

	ctx := context.Background()
	got, err := federated.PublicKeys(ctx)
	if err != nil {
		t.Fatalf("PublicKeys: %v", err)
	}

	var decoded stoke.JWKSet
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(decoded.Keys) != 2 {
		t.Fatalf("expected 2 keys (local + peer), got %d", len(decoded.Keys))
	}
	kids := make(map[string]bool)
	for _, k := range decoded.Keys {
		kids[k.KeyId] = true
	}
	if !kids["p-0"] || !kids["p-1"] {
		t.Errorf("expected keys p-0 and p-1, got %v", kids)
	}
}

// TestFederatedTokenIssuer_ParseClaims_VerifiesTokenFromPeer ensures that a token signed by
// a "peer" replica (key B) is verified by the federated issuer's merged JWKS (local key A + peer key B).
func TestFederatedTokenIssuer_ParseClaims_VerifiesTokenFromPeer(t *testing.T) {
	pubB, privB, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	jwkB := stoke.CreateJWK().FromPublicKey(pubB)
	jwkB.KeyId = "p-peer"
	peerJWKS := stoke.JWKSet{
		Expires: time.Now().Add(time.Hour),
		Keys:    []*stoke.JWK{jwkB},
	}
	peerBytes, err := json.Marshal(peerJWKS)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pkeys" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(peerBytes)
	}))
	defer srv.Close()

	// Local "replica" has only key A (mock); federated merge will add peer key B.
	localJWKS := stoke.JWKSet{
		Expires: time.Now().Add(time.Hour),
		Keys:    []*stoke.JWK{{KeyId: "p-0", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x0", Y: "y0"}},
	}
	localBytes, _ := json.Marshal(localJWKS)
	inner := &mockFederatedInner{publicKeysBytes: localBytes}
	discoverer := &cluster.StaticDiscoverer{URLs: []string{srv.URL}}
	federated := NewFederatedTokenIssuer(inner, discoverer, nil, "", 0)

	// Sign a token as "peer" would (with key B).
	claims := &stoke.Claims{}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
	claims.Issuer = "test"
	claims.StokeClaims = map[string]string{"kid": "p-peer"}
	token, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(privB)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	// Verify with federated issuer (merged set includes B).
	ctx := context.Background()
	reqClaims := &stoke.Claims{}
	parsed, err := federated.ParseClaims(ctx, token, reqClaims, jwt.WithExpirationRequired())
	if err != nil {
		t.Fatalf("ParseClaims (token from peer): %v", err)
	}
	if !parsed.Valid {
		t.Error("parsed token should be valid")
	}
}
