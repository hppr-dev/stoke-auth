package key

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
)

type EdDSAKeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey ed25519.PublicKey
	KeyMeta
}

func (k *EdDSAKeyPair) Generate() error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	k.PrivateKey = priv
	k.PublicKey = pub
	return err
}

func (k *EdDSAKeyPair) PublicString() string {
	return base64.StdEncoding.EncodeToString(k.PublicKey)
}

func (k *EdDSAKeyPair) Encode() string {
	return base64.StdEncoding.EncodeToString(k.PrivateKey)
}

func (k *EdDSAKeyPair) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}
	k.PrivateKey = b
	k.PublicKey = k.PrivateKey.Public().(ed25519.PublicKey)

	return nil
}

func (k *EdDSAKeyPair) Key() ed25519.PrivateKey {
	return k.PrivateKey
}

func (k *EdDSAKeyPair) SigningMethod() jwt.SigningMethod {
	return jwt.GetSigningMethod("EdDSA")
}
