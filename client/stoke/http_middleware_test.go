package stoke_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"hppr.dev/stoke"
)

type testHandler struct {
	called bool
}

func (h *testHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	h.called = true
}

func TestAuthHandlerSyntax(t *testing.T) {
	store := stoke.NewTestPublicKeyStore(nil)
	stoke.Auth(&testHandler{}, store, stoke.RequireToken())
}

func TestAuthFuncHandlerSyntax(t *testing.T) {
	store := stoke.NewTestPublicKeyStore(nil)
	stoke.AuthFunc(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request){}), store, stoke.RequireToken())
}

func TestInnerHandlerIsCalledWithValidToken(t *testing.T) {
	inner := &testHandler{}
	store := stoke.NewTestPublicKeyStore(nil)
	outer := stoke.Auth(inner, store, stoke.RequireToken())

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Add("Authorization", "Bearer invalid")
	res := httptest.NewRecorder()

	outer.ServeHTTP(res, req)

	if !inner.called {
		t.Fatal("Inner handler not called")
	}

	if res.Result().StatusCode != 200 {
		t.Fatalf("Return code was not 200: %v", res.Result())
	}
}

func TestInnerHandlerIsNotCalledWithInvalidToken(t *testing.T) {
	inner := &testHandler{}
	store := stoke.NewTestPublicKeyStore(nil)
	outer := stoke.Auth(inner, store, stoke.RequireToken())
	store.SetInvalid()

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Add("Authorization", "Bearer invalid")
	res := httptest.NewRecorder()

	outer.ServeHTTP(res, req)

	if inner.called {
		t.Fatal("Inner handler called")
	}

	if res.Result().StatusCode != 401 {
		t.Fatalf("Return code was not 401: %v", res.Result())
	}
}

func TestInnerHandlerIsNotCalledWithNoToken(t *testing.T) {
	inner := &testHandler{}
	store := stoke.NewTestPublicKeyStore(nil)
	outer := stoke.Auth(inner, store, stoke.RequireToken())
	store.SetInvalid()

	req := httptest.NewRequest("GET", "http://localhost", nil)
	res := httptest.NewRecorder()

	outer.ServeHTTP(res, req)

	if inner.called {
		t.Fatal("Inner handler called")
	}

	if res.Result().StatusCode != 422 {
		t.Fatalf("Return code was not 422: %v", res.Result())
	}
}
