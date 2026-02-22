package cluster

import (
	"encoding/json"
	"net/http"
	"strings"

	"hppr.dev/stoke"
)

// MergeJWKS parses localJWKS as a JWKSet, fetches JWKS from each peer at
// peerURL + "/api/pkeys", merges all keys deduplicating by KeyId, and returns
// the combined JWKSet as JSON. Expires is set to the earliest expiry among
// local and all fetched sets. If httpClient is nil, http.DefaultClient is used.
// Peer fetch failures (non-200 or decode error) cause that peer to be skipped;
// the merge still succeeds.
func MergeJWKS(localJWKS []byte, peerURLs []string, httpClient *http.Client) ([]byte, error) {
	client := httpClient
	if client == nil {
		client = http.DefaultClient
	}

	var local stoke.JWKSet
	if err := json.Unmarshal(localJWKS, &local); err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var keys []*stoke.JWK
	expires := local.Expires

	for _, k := range local.Keys {
		if k == nil {
			continue
		}
		if seen[k.KeyId] {
			continue
		}
		seen[k.KeyId] = true
		keys = append(keys, k)
	}

	for _, baseURL := range peerURLs {
		u := strings.TrimSuffix(baseURL, "/") + "/api/pkeys"
		resp, err := client.Get(u)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			continue
		}
		var peer stoke.JWKSet
		if err := json.NewDecoder(resp.Body).Decode(&peer); err != nil {
			_ = resp.Body.Close()
			continue
		}
		_ = resp.Body.Close()

		if !peer.Expires.IsZero() && (expires.IsZero() || peer.Expires.Before(expires)) {
			expires = peer.Expires
		}
		for _, k := range peer.Keys {
			if k == nil {
				continue
			}
			if seen[k.KeyId] {
				continue
			}
			seen[k.KeyId] = true
			keys = append(keys, k)
		}
	}

	out := stoke.JWKSet{
		Expires: expires,
		Keys:    keys,
	}
	return json.Marshal(out)
}
