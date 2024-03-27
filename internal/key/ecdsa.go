package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type ECDSAKeyPair struct {
	NumBits int
	PrivateKey *ecdsa.PrivateKey
	PublicKey *ecdsa.PublicKey
	KeyMeta
}

func (k *ECDSAKeyPair) Generate() error {
	log.Println("Generating new ECDSA keypair...")
	priv, err := ecdsa.GenerateKey(k.getCurve(), rand.Reader)
	k.PrivateKey = priv
	k.PublicKey = &priv.PublicKey
	return err
}

func (k *ECDSAKeyPair) PublicString() string {
	s, _ := x509.MarshalPKIXPublicKey(k.PublicKey)
	return base64.StdEncoding.EncodeToString(s)
}

func (k *ECDSAKeyPair) Encode() string {
	s, _ := x509.MarshalECPrivateKey(k.PrivateKey)
	return base64.StdEncoding.EncodeToString(s)
}

func (k *ECDSAKeyPair) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	k.PrivateKey, err = x509.ParseECPrivateKey(b)
	if err != nil {
		return err
	}

	k.PublicKey = &k.PrivateKey.PublicKey

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
	log.Println("Number of bits not set to 256, 384, or 512. Using default 256.")
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
	log.Println("Number of bits not set to 256, 384, or 512. Using default 256.")
	return elliptic.P256()
}

