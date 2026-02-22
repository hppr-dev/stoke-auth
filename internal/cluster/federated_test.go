package cluster

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hppr.dev/stoke"
)

func TestMergeJWKS_LocalOnly(t *testing.T) {
	exp := time.Now().Add(time.Hour)
	local := stoke.JWKSet{
		Expires: exp,
		Keys: []*stoke.JWK{
			{KeyId: "p-0", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x", Y: "y"},
		},
	}
	localBytes, err := json.Marshal(local)
	if err != nil {
		t.Fatal(err)
	}

	// peerURLs nil
	got, err := MergeJWKS(localBytes, nil, nil)
	if err != nil {
		t.Fatalf("MergeJWKS(nil peers): %v", err)
	}
	var decoded stoke.JWKSet
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(decoded.Keys) != 1 || decoded.Keys[0].KeyId != "p-0" {
		t.Errorf("expected one key p-0, got %d keys", len(decoded.Keys))
		if len(decoded.Keys) > 0 {
			t.Errorf("first key kid: %q", decoded.Keys[0].KeyId)
		}
	}
	if !decoded.Expires.Equal(exp) {
		t.Errorf("expires: got %v, want %v", decoded.Expires, exp)
	}

	// peerURLs empty
	got2, err := MergeJWKS(localBytes, []string{}, nil)
	if err != nil {
		t.Fatalf("MergeJWKS(empty peers): %v", err)
	}
	var decoded2 stoke.JWKSet
	if err := json.Unmarshal(got2, &decoded2); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(decoded2.Keys) != 1 || decoded2.Keys[0].KeyId != "p-0" {
		t.Errorf("empty peers: expected one key p-0, got %d keys", len(decoded2.Keys))
	}
}

func TestMergeJWKS_OnePeer(t *testing.T) {
	expLocal := time.Now().Add(2 * time.Hour)
	local := stoke.JWKSet{
		Expires: expLocal,
		Keys: []*stoke.JWK{
			{KeyId: "p-0", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x0", Y: "y0"},
		},
	}
	localBytes, err := json.Marshal(local)
	if err != nil {
		t.Fatal(err)
	}

	expPeer := time.Now().Add(1 * time.Hour)
	peerSet := stoke.JWKSet{
		Expires: expPeer,
		Keys: []*stoke.JWK{
			{KeyId: "p-1", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x1", Y: "y1"},
		},
	}
	peerBytes, err := json.Marshal(peerSet)
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

	got, err := MergeJWKS(localBytes, []string{srv.URL}, nil)
	if err != nil {
		t.Fatalf("MergeJWKS: %v", err)
	}

	var decoded stoke.JWKSet
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(decoded.Keys) != 2 {
		t.Fatalf("expected 2 keys (p-0, p-1), got %d", len(decoded.Keys))
	}
	kids := make(map[string]bool)
	for _, k := range decoded.Keys {
		kids[k.KeyId] = true
	}
	if !kids["p-0"] || !kids["p-1"] {
		t.Errorf("expected keys p-0 and p-1, got %v", kids)
	}
	// Earliest expiry should be peer's (1h < 2h)
	if !decoded.Expires.Equal(expPeer) {
		t.Errorf("expires: got %v, want earliest %v", decoded.Expires, expPeer)
	}
}

func TestMergeJWKS_DedupByKeyId(t *testing.T) {
	exp := time.Now().Add(time.Hour)
	local := stoke.JWKSet{
		Expires: exp,
		Keys: []*stoke.JWK{
			{KeyId: "p-0", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x0", Y: "y0"},
		},
	}
	localBytes, _ := json.Marshal(local)

	// Peer returns same kid "p-0" (duplicate)
	peerSet := stoke.JWKSet{
		Expires: exp,
		Keys: []*stoke.JWK{
			{KeyId: "p-0", KeyType: "EC", Use: "sig", Curve: "P-256", X: "other", Y: "other"},
			{KeyId: "p-1", KeyType: "EC", Use: "sig", Curve: "P-256", X: "x1", Y: "y1"},
		},
	}
	peerBytes, _ := json.Marshal(peerSet)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(peerBytes)
	}))
	defer srv.Close()

	got, err := MergeJWKS(localBytes, []string{srv.URL}, nil)
	if err != nil {
		t.Fatalf("MergeJWKS: %v", err)
	}
	var decoded stoke.JWKSet
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// p-0 once (from local), p-1 from peer
	if len(decoded.Keys) != 2 {
		t.Fatalf("dedup by kid: expected 2 keys, got %d", len(decoded.Keys))
	}
	kids := make(map[string]bool)
	for _, k := range decoded.Keys {
		kids[k.KeyId] = true
	}
	if !kids["p-0"] || !kids["p-1"] {
		t.Errorf("expected p-0 and p-1, got %v", kids)
	}
}
