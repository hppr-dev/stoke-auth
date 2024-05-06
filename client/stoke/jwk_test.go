package stoke_test

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"hppr.dev/stoke"
)

func TestRSAJWK(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}

	jwk := stoke.CreateJWK().FromPublicKey(rsaKey.Public())

	postPubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("Could not convert JWK to public key: %v", err)
	}

	if rsaPubKey, ok := postPubKey.(*rsa.PublicKey); !ok || !rsaPubKey.Equal(&rsaKey.PublicKey) {
		t.Fatalf("Converted key did not match: %v %T, %v", ok, postPubKey, postPubKey)
	}
}

func TestECDSA256JWK(t *testing.T) {
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}

	jwk := stoke.CreateJWK().FromPublicKey(ecdsaKey.Public())

	postPubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("Could not convert JWK to public key: %v", err)
	}

	if ecdsaPubKey, ok := postPubKey.(*ecdsa.PublicKey); !ok || !ecdsaPubKey.Equal(&ecdsaKey.PublicKey) {
		t.Fatalf("Converted key did not match: %v %T, %v", ok, postPubKey, postPubKey)
	}
}

func TestECDSA384JWK(t *testing.T) {
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}

	jwk := stoke.CreateJWK().FromPublicKey(ecdsaKey.Public())

	postPubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("Could not convert JWK to public key: %v", err)
	}

	if ecdsaPubKey, ok := postPubKey.(*ecdsa.PublicKey); !ok || !ecdsaPubKey.Equal(&ecdsaKey.PublicKey) {
		t.Fatalf("Converted key did not match: %v %T, %v", ok, postPubKey, postPubKey)
	}
}

func TestECDSA521JWK(t *testing.T) {
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}

	jwk := stoke.CreateJWK().FromPublicKey(ecdsaKey.Public())

	postPubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("Could not convert JWK to public key: %v", err)
	}

	if ecdsaPubKey, ok := postPubKey.(*ecdsa.PublicKey); !ok || !ecdsaPubKey.Equal(&ecdsaKey.PublicKey) {
		t.Fatalf("Converted key did not match: %v %T, %v", ok, postPubKey, postPubKey)
	}
}

func TestEdDSAJWK(t *testing.T) {
	preEddsaKey, eddsaKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Could not generate key: %v", err)
	}

	jwk := stoke.CreateJWK().FromPublicKey(eddsaKey.Public())

	postPubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("Could not convert JWK to public key: %v", err)
	}

	if eddsaPubKey, ok := postPubKey.(ed25519.PublicKey); !ok || !eddsaPubKey.Equal(preEddsaKey) {
		t.Fatalf("Converted key did not match: %v %T, %v", ok, postPubKey, postPubKey)
	}
}

func TestNilJWK(t *testing.T) {
	jwk := stoke.CreateJWK()
	// dsa is a legacy algorithm
	if pubKey := jwk.FromPublicKey(dsa.PublicKey{}); pubKey != nil {
		t.Fatalf("Returned non-nil pubkey: %T %v", pubKey, pubKey)
	}
}

func TestBadKeyType(t *testing.T) {
	jwk := stoke.JWK{
		KeyType: "BAD",
	}

	if _, err := jwk.ToPublicKey(); err == nil {
		t.Fatal("Did not return an error")
	}
}

func TestBadEncoding(t *testing.T) {
	goodStr :=  "ATX6CAe3lTbvBMLjftp3sC9BDLP7AZvAaR9RjYyDbYUiTvgxmEvejyAUPXJOmsFpQ9zstrJ-YcH1VnvbWlW0f9zn"
	badStr := "BAD/ENCODE"
	badX := stoke.JWK{ KeyType: "EC", Curve :"P-256", X : badStr, Y : goodStr }
	if _, err := badX.ToPublicKey(); err == nil {
		t.Log("Bad X encoding did not return an error")
		t.Fail()
	}

	badY := stoke.JWK{ KeyType: "EC", Curve :"P-256", Y : badStr, X : goodStr }
	if _, err := badY.ToPublicKey(); err == nil {
		t.Log("Bad Y encoding did not return an error")
		t.Fail()
	}

	badN := stoke.JWK{ KeyType: "RSA", N : badStr, E : goodStr }
	if _, err := badN.ToPublicKey(); err == nil {
		t.Log("Bad N encoding did not return an error")
		t.Fail()
	}

	badE := stoke.JWK{ KeyType: "RSA", N : goodStr, E : badStr }
	if _, err := badE.ToPublicKey(); err == nil {
		t.Log("Bad E encoding did not return an error")
		t.Fail()
	}
}

func TestBadCurve(t *testing.T) {
	goodStr :=  "ATX6CAe3lTbvBMLjftp3sC9BDLP7AZvAaR9RjYyDbYUiTvgxmEvejyAUPXJOmsFpQ9zstrJ-YcH1VnvbWlW0f9zn"
	badStr := "BAD/ENCODE"

	badCurve := stoke.JWK{ KeyType: "EC", Curve :"bad", X : badStr, Y : goodStr }
	if _, err := badCurve.ToPublicKey(); err == nil {
		t.Log("Bad curve did not return an error")
		t.Fail()
	}
}
