package key

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type ECDSAKeyPair struct {
	NumBits int
	PrivateKey *ecdsa.PrivateKey
	KeyMeta
	Logger zerolog.Logger
}

func (k *ECDSAKeyPair) Generate() (KeyPair[*ecdsa.PrivateKey], error) {
	k.Logger.Info().Msg("Generating new ECDSA keypair...")

	priv, err := ecdsa.GenerateKey(k.getCurve(), rand.Reader)
	return &ECDSAKeyPair{
		NumBits: k.NumBits,
		PrivateKey: priv,
	}, err
}

func (k *ECDSAKeyPair) PublicString() string {
	s, _ := x509.MarshalPKIXPublicKey(&k.PrivateKey.PublicKey)
	return base64.StdEncoding.EncodeToString(s)
}

func (k *ECDSAKeyPair) PublicKey() crypto.PublicKey {
	return &k.PrivateKey.PublicKey
}

func (k *ECDSAKeyPair) Encode() string {
	s, _ := x509.MarshalECPrivateKey(k.PrivateKey)
	return base64.StdEncoding.EncodeToString(s)
}

func (k *ECDSAKeyPair) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		k.Logger.Error().Err(err).Msg("Decoding base64 private key failed")
		return err
	}

	k.PrivateKey, err = x509.ParseECPrivateKey(b)
	if err != nil {
		k.Logger.Error().Err(err).Msg("Decoding ECDSA private key failed")
		return err
	}

	return nil
}

func (k *ECDSAKeyPair) Key() *ecdsa.PrivateKey {
	return k.PrivateKey
}

func (k *ECDSAKeyPair) SigningMethod() jwt.SigningMethod {
	switch k.NumBits {
	case 256:
		return jwt.GetSigningMethod("ES256")
	case 384:
		return jwt.GetSigningMethod("ES384")
	case 512:
		return jwt.GetSigningMethod("ES512")
	}

	k.Logger.Info().Msg("Number of bits not set to 256, 384, or 512. Using default 256.")
	return jwt.GetSigningMethod("ES256")
}


func (k *ECDSAKeyPair) getCurve() elliptic.Curve {
	switch k.NumBits {
	case 256:
		return elliptic.P256()
	case 384:
		return elliptic.P384()
	case 512:
		return elliptic.P521()
	}

	k.Logger.Info().Msg("Number of bits not set to 256, 384, or 512. Using default 256.")
	return elliptic.P256()
}

