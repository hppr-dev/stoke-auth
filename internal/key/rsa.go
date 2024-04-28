package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type RSAKeyPair struct {
	NumBits int
	PrivateKey *rsa.PrivateKey
	KeyMeta
	Logger zerolog.Logger
}

func (k *RSAKeyPair) Generate() (KeyPair[*rsa.PrivateKey], error) {
	k.Logger.Info().Msg("Generating RSA key...")

	if k.NumBits != 256 && k.NumBits != 384 && k.NumBits != 512 {
		k.Logger.Warn().Msg("Number of bits not set to 256, 384, or 512. Setting to default 256.")
		k.NumBits = 256
	}

	priv, err := rsa.GenerateKey(rand.Reader, k.NumBits)
	return &RSAKeyPair{
		NumBits: k.NumBits,
		PrivateKey: priv,
	}, err
}

func (k *RSAKeyPair) PublicString() string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&k.PrivateKey.PublicKey))
}

func (k *RSAKeyPair) PublicKey() crypto.PublicKey {
	return &k.PrivateKey.PublicKey
}

func (k *RSAKeyPair) Encode() string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(k.PrivateKey))
}

func (k *RSAKeyPair) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		k.Logger.Error().Err(err).Msg("Error decoding base64 RSA private key")
		return err
	}

	k.PrivateKey, err = x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		k.Logger.Error().Err(err).Msg("Error decoding PKCS1 RSA private key")
		return err
	}

	return nil
}

func (k *RSAKeyPair) Key() *rsa.PrivateKey {
	return k.PrivateKey
}

func (k *RSAKeyPair) SigningMethod() jwt.SigningMethod {
	switch k.NumBits {
	case 256:
		return jwt.GetSigningMethod("PS256")
	case 384:
		return jwt.GetSigningMethod("PS384")
	case 512:
		return jwt.GetSigningMethod("PS512")
	}
	k.Logger.Info().Msg("Number of bits not set to 256, 384, or 512. Using default 256.")
	return jwt.GetSigningMethod("PS256")
}
