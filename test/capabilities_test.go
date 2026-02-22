package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const capabilitiesURL = "http://localhost:8080/api/capabilities"
const availableProvidersURL = "http://localhost:8080/api/available_providers"
const loginURL = "http://localhost:8080/api/login"

func TestCapabilities_ReturnsCapabilities(t *testing.T) {
	token, err := loginForToken(t)
	if err != nil {
		t.Skipf("skipping capabilities test: could not obtain token (is server running?): %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, capabilitiesURL, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("capabilities request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("capabilities: got status %d, want 200", resp.StatusCode)
	}

	var result struct {
		Capabilities []string `json:"capabilities"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode capabilities response: %v", err)
	}

	if result.Capabilities == nil {
		t.Error("response must include capabilities array")
	}
}

func loginForToken(t *testing.T) (string, error) {
	t.Helper()
	loginBody := []byte(`{"username":"tester","password":"tester","provider":"","required_claims":[{"stk":""}],"filter_claims":["stk"]}`)
	resp, err := http.Post(loginURL, "application/json", bytes.NewReader(loginBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login returned %d", resp.StatusCode)
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Token, nil
}

func TestAvailableProviders_ReturnsProvidersAndBaseAdminPath(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, availableProvidersURL, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("skipping available_providers test: could not reach server (is server running?): %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("available_providers: got status %d, want 200", resp.StatusCode)
	}

	var result struct {
		Providers     []struct {
			Name         string `json:"name"`
			ProviderType string `json:"provider_type"`
			TypeSpec     string `json:"type_spec"`
		} `json:"providers"`
		BaseAdminPath string `json:"base_admin_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode available_providers response: %v", err)
	}

	if result.Providers == nil {
		t.Error("response must include providers array")
	}
	// base_admin_path is optional; when present it must be a string (already decoded).
	// When server is configured with base_admin_path (e.g. /auth), it will be present.
	_ = result.BaseAdminPath
}
