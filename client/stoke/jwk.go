package stoke

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"
)

// Customized Json Web Key (JWK) Set.
// {
//   "exp" : "next update", // When to pull the next certificates
//   "keys": [
//      { "kty": "RSA",", "n": "BASE64_N", "e": "BASE64_E", "use" : "sig", "kid" : "keyID", "x5c" : [ "PKIX_CERT" ] }],      // EXAMPLE RSA
//      { "kty": "EC", "crv": "P-256(CURVE_TYPE)", "x" : "BASE64URL_X", "y": "BASE64URL_Y" "use" : "sig", "kid" : "keyID" }] // EXAMPLE EC
//   ]
// }
// See also: https://datatracker.ietf.org/doc/html/rfc7517#autoid-5
// See also: https://www.iana.org/assignments/jose/jose.xhtml#web-key-types
type JWKSet struct {
	Expires time.Time `json:"exp"`
	Keys    []*JWK    `json:"keys"`
}

// Represents a Json Web Key that conforms to the proposed standard in RF7517
// See also: https://datatracker.ietf.org/doc/html/rfc7517#autoid-5
type JWK struct {
	KeyType   string `json:"kty,omitempty"`
	Use       string `json:"use,omitempty"`
	KeyId     string `json:"kid,omitempty"`
	// EC fields
	Curve     string `json:"crv,omitempty"`
	X         string `json:"x,omitempty"`
	Y         string `json:"y,omitempty"`
	// RSA fields
	N         string `json:"n,omitempty"`
	E         string `json:"e,omitempty"`
	// OKP (ed25519) uses Curve and X
}

// Converts a JWK into a crypto.PublicKey
func (j *JWK) ToPublicKey() (crypto.PublicKey, error) {
	switch j.KeyType {
	case "EC":
		return j.ToECDSA()
	case "RSA":
		return j.ToRSA()
	case "OKP":
		return j.ToEdDSA()
	default:
		return nil, fmt.Errorf("Unknown key type: %s", j.KeyType)
	}
}

// Creates a JWK.
// Allows for method chaining: CreateJWK().FromECDSA(key), CreateJWK().FromEdDSA(key), CreateJWK().FromRSA(key)
func CreateJWK() *JWK {
	return &JWK{}
}

// Loads a JWK with info from a crypto.PublicKey
// Supported key types are: *ecdsa.PublicKey, *rsa.PublicKey, and ed25519.PublicKey
func (j *JWK) FromPublicKey(key crypto.PublicKey) *JWK {
	switch key.(type) {
	case *ecdsa.PublicKey:
		return j.FromECDSA(key.(*ecdsa.PublicKey))
	case *rsa.PublicKey:
		return j.FromRSA(key.(*rsa.PublicKey))
	case ed25519.PublicKey:
		return j.FromEdDSA(key.(ed25519.PublicKey))
	}
	return nil
}

// Fills a JWK with an ecdsa key
func (j *JWK) FromECDSA(key *ecdsa.PublicKey) *JWK {
	j.KeyType = "EC"
	j.Use = "sig"
	j.Curve = key.Curve.Params().Name
	j.X = base64.URLEncoding.EncodeToString(key.X.Bytes())
	j.Y = base64.URLEncoding.EncodeToString(key.Y.Bytes())
	return j
}

// Converts JWK to an ECDSA key
func (j *JWK) ToECDSA() (*ecdsa.PublicKey, error) {
	var curve elliptic.Curve
	switch j.Curve {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("Unknown Curve type: %s", j.Curve)
	}

	xBig, err := stringToBigInt(j.X)
	if err != nil {
		return nil, err
	}

	yBig,err := stringToBigInt(j.Y)
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     xBig,
		Y:     yBig,
	}, nil
}

// Fills a JWK with an EdDSA key
func (j *JWK) FromEdDSA(key ed25519.PublicKey) *JWK {
	j.KeyType = "EC"
	j.Use = "sig"
	j.Curve = "ed25519"
	j.X = base64.URLEncoding.EncodeToString(key)

	return j
}

// Converts JWK to an EdDSA (ed25519) key
func (j *JWK) ToEdDSA() (ed25519.PublicKey, error) {
	decoded, err := base64.URLEncoding.DecodeString(j.X)
	return ed25519.PublicKey(decoded), err
}

// Fills a JWK with an RSA key
func (j *JWK) FromRSA(key *rsa.PublicKey) *JWK {
	j.KeyType = "RSA"
	j.Use = "sig"
	j.N = base64.URLEncoding.EncodeToString(key.N.Bytes())

	exponentBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(exponentBytes, uint32(key.E))
	j.E = base64.URLEncoding.EncodeToString(exponentBytes)

	return j
}

// Converts JWK to an RSA key
func (j *JWK) ToRSA() (*rsa.PublicKey, error) {
	decodedE, err := base64.URLEncoding.DecodeString(j.E)
	if err != nil {
		return nil, err
	}

	bigN, err := stringToBigInt(j.N)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: bigN,
		E: int(binary.BigEndian.Uint32(decodedE)),
	}, nil
}


// Converts a given string in base64 URL encoding to a *big.Int
func stringToBigInt(b string) (*big.Int, error) {
	convB, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(convB), nil
}
