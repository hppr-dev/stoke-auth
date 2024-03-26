package key

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)


type PrivateKey interface {
	*rsa.PrivateKey|*ecdsa.PrivateKey|ed25519.PrivateKey
}

type PublicKey interface {
	*rsa.PublicKey|*ecdsa.PublicKey|ed25519.PublicKey
}

type KeyPair[P PrivateKey, K PublicKey] interface {
	Generate() error
	PublicString() string
	Encode() string
	Decode(string) error
	Keys() (P, K)
	SigningMethod() jwt.SigningMethod
}
