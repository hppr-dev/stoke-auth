package key

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PrivateKey interface {
	*rsa.PrivateKey|*ecdsa.PrivateKey|ed25519.PrivateKey
}

type KeyPair[P PrivateKey] interface {
	Generate() error
	PublicString() string
	Encode() string
	Decode(string) error
	Key() P
	PublicKey() crypto.PublicKey
	SigningMethod() jwt.SigningMethod
	SetExpires(time.Time)
	ExpiresAt() time.Time
	SetRenews(time.Time)
	RenewsAt() time.Time
}
