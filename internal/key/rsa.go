package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type RSAKeyPair struct {
	NumBits int
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

func (k *RSAKeyPair) Generate() error {
	if k.NumBits != 256 && k.NumBits != 384 && k.NumBits != 512 {
		log.Println("Number of bits not set to 256, 384, or 512. Setting to default 256.")
		k.NumBits = 256
	}
	priv, err := rsa.GenerateKey(rand.Reader, k.NumBits)
	k.PrivateKey = priv
	k.PublicKey = &priv.PublicKey
	return err
}

func (k *RSAKeyPair) PublicString() string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(k.PublicKey))
}

func (k *RSAKeyPair) Encode() string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(k.PrivateKey))
}

func (k *RSAKeyPair) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}
	k.PrivateKey, err = x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return err
	}

	k.PublicKey = &k.PrivateKey.PublicKey

	return nil
}

func (k *RSAKeyPair) Keys() (*rsa.PrivateKey, *rsa.PublicKey) {
	return k.PrivateKey, k.PublicKey
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
	log.Println("Number of bits not set to 256, 384, or 512. Using default 256.")
	return jwt.GetSigningMethod("PS256")
}
