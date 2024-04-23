package stoke

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestPublicKeyStore is a mocked out public keystore
// See https://jwt.io#debugger to generate a DefaultTokenStr
// USE IN DEVELOPMENT OR FOR TESTING ONLY
type TestPublicKeyStore struct {
	DefaultTokenStr string
	reject bool
	invalid bool
}

// Treat all calls to ParseClaims as an error
func (t *TestPublicKeyStore) SetReject() {
	t.reject = true
}

// Allow ParseClaims to return a token
func (t *TestPublicKeyStore) SetAllow() {
	t.reject = false
}

// Treat all calls to ParseClaims as valid tokens.
func (t *TestPublicKeyStore) SetValid() {
	t.invalid = false
}

// Treat all calls to ParseClaims as invalid tokens.
func (t *TestPublicKeyStore) SetInvalid() {
	t.invalid = true
}

// NOOP
func (TestPublicKeyStore) Init(ctx context.Context) error { return nil }

// Turns off parsing/validating claims, call SetReject to deny all access.
// If tokenStr is empty, DefaultTokenStr will be used.
// tokenStr must be a parsable JWT token or an empty string.
//
// TestPublicKeyStore implements PublicKeyStore
func (t TestPublicKeyStore) ParseClaims(ctx context.Context, tokenStr string, claims *Claims, _ ...jwt.ParserOption) (*jwt.Token, error) {
	if t.reject {
		return nil, fmt.Errorf("Reject set.")
	}

	if tokenStr == "" {
		tokenStr = t.DefaultTokenStr
	}

	token, _ := jwt.ParseWithClaims(tokenStr, claims.New(), fakeKeyFunc,
		jwt.WithoutClaimsValidation(),
	)
	token.Valid = !t.invalid
	return token, nil
}

// NOOP
func fakeKeyFunc(*jwt.Token) (interface{}, error) { return nil, nil }
