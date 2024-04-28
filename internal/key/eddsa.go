package key

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type EdDSAKeyPair struct {
	PrivateKey ed25519.PrivateKey
	KeyMeta
	Logger zerolog.Logger
}

func (k *EdDSAKeyPair) Generate() (KeyPair[ed25519.PrivateKey], error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	return &EdDSAKeyPair{
		PrivateKey: priv,
	}, err
}

func (k *EdDSAKeyPair) PublicString() string {
	b := k.PrivateKey.Public().(ed25519.PublicKey)
	return base64.StdEncoding.EncodeToString(b)
}

func (k *EdDSAKeyPair) PublicKey() crypto.PublicKey {
	return k.PrivateKey.Public()
}

func (k *EdDSAKeyPair) Encode() string {
	return base64.StdEncoding.EncodeToString(k.PrivateKey)
}

func (k *EdDSAKeyPair) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		k.Logger.Error().Err(err).Msg("Could not decode EdDSA private key")
		return err
	}
	k.PrivateKey = b

	return nil
}

func (k *EdDSAKeyPair) Key() ed25519.PrivateKey {
	return k.PrivateKey
}

func (k *EdDSAKeyPair) SigningMethod() jwt.SigningMethod {
	return jwt.GetSigningMethod("EdDSA")
}
