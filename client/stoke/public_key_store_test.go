package stoke_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"hppr.dev/stoke"
)

var token = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdGsiLCJleHAiOjQ4NzAzNjQ2NDksImp0aSI6ImsyIiwiZSI6InNhZG1pbkBsb2NhbGhvc3QiLCJraWQiOiIwIiwic3JvbCI6InNwciIsInUiOiJzYWRtaW4iLCJuIjoiU3Rva2UgQWRtaW4ifQ.YKKq9y2A4bZJ3WK_ntDqjQE4THPkY7RRjKR6htMIgbgbq7Et5G-5Ba5QwwvtaF2JtB36f5YRlwpEMXV_DhSZEQ"

type testServer struct {}

func (testServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte(`
	{
		"exp":"2266-08-11T05:50:55-05:00",
		"keys":[
			{
				"kty":"EC",
				"use":"sig",
				"kid":"p-0",
				"crv":"P-256",
				"x":"Ja9L9-ew9h-NZSGzCN3QSbzH3gg96Grl0wh-4IH5F7U=",
				"y":"6A8YccjPtbVD-jqQTTTQlSgFJHU60Xphgbqs65vZ5is="
			},
			{
				"kty":"EC",
				"use":"sig",
				"kid":"p-1",
				"crv":"P-256",
				"x":"_TPMCEa_V2qZIg6UKNRDGXz-Pk1WZwzcPQzX0qhClo0=",
				"y":"y4rTZtP_LpZ3ocpOAJ5yxHRuoGLprEe67gm8NR6f9zQ="
			}
		]
	}`))
}

func TestPerRequestParseClaims(t *testing.T) {
	server := httptest.NewServer(testServer{})
	store, err := stoke.NewPerRequestPublicKeyStore(server.URL, context.Background())
	if err != nil {
		t.Fatalf("Could not create new store: %v", err)
	}

	token, err := store.ParseClaims(context.Background(), token, stoke.RequireToken().WithClaim("srol", "spr"))
	if err != nil {
		t.Fatalf("Parsing claims returned an error: %v", err)
	}

	if claims, ok := token.Claims.(*stoke.Claims); !ok || claims.StokeClaims["n"] != "Stoke Admin" || claims.StokeClaims["srol"] != "spr" || claims.StokeClaims["u"] != "sadmin" {
		t.Fatalf("Claims did not match: %v", claims)
	}
}

func TestWebRequestParseClaims(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := httptest.NewServer(testServer{})
	store, err := stoke.NewWebCachePublicKeyStore(server.URL, ctx)
	if err != nil {
		t.Fatalf("Could not create new store: %v", err)
	}

	token, err := store.ParseClaims(context.Background(), token, stoke.RequireToken().WithClaim("srol", "spr"))
	if err != nil {
		t.Fatalf("Parsing claims returned an error: %v", err)
	}

	if claims, ok := token.Claims.(*stoke.Claims); !ok || claims.StokeClaims["n"] != "Stoke Admin" || claims.StokeClaims["srol"] != "spr" || claims.StokeClaims["u"] != "sadmin" {
		t.Fatalf("Claims did not match: %v", claims)
	}
}
