package key_test

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"stoke/internal/key"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateHappy(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	newEd, err := ed.Generate()
	if err != nil || newEd == nil {
		t.Fatalf("Failed to generate EdDSAKeyPair: %v", err)
	}

	newEc, err := ec.Generate()
	if err != nil || newEc == nil {
		t.Fatalf("Failed to generate ECDSAKeyPair: %v", err)
	}

	ec.NumBits = 384
	newEc, err = ec.Generate()
	if err != nil || newEc == nil {
		t.Fatalf("Failed to generate ECDSAKeyPair: %v", err)
	}

	ec.NumBits = 512
	newEc, err = ec.Generate()
	if err != nil || newEc == nil {
		t.Fatalf("Failed to generate ECDSAKeyPair: %v", err)
	}

	// Bad num bits defaults to 256
	ec.NumBits = 1
	newEc, err = ec.Generate()
	if err != nil || newEc == nil {
		t.Fatalf("Failed to generate ECDSAKeyPair: %v", err)
	}

	if newEc.Key().Params().BitSize != 256 {
		t.Fatalf("Defaultd bit sized ECDSAKeyPair was not 256: %d", newEc.Key().Params().BitSize)
	}

	newRs, err := rs.Generate()
	if err != nil || newRs == nil {
		t.Fatalf("Failed to generate RSSAKeyPair: %v", err)
	}

	// Bad num bits defaults to 256
	rs.NumBits = 1
	newRs, err = rs.Generate()
	if err != nil || newRs == nil {
		t.Fatalf("Failed to generate RSSAKeyPair: %v", err)
	}

	if newRs.Key().Size() != (256 / 8) {
		t.Fatalf("RSA was not what was expected: %d", newRs.Key().Size())
	}
}

func TestPublicStringHappy(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	expEd := "tftbPsVuL85T0jxjdzkxErXQSmlt0h4zERiS+a5OiAA="
	expEc := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEbJD2pRVSiKjnrnkEY+Cs+0n3RcJjfhO5hzz8H3ilMgvshY0xdkBEuEwIpIJFYXUCladiye+hmWOgFxOen8EoTw=="
	expRs := "MCgCIQC+9cu5fmrUYhkbUNoEP5S+sUqv6KnaB32IBa3Y9HlsdQIDAQAB"

	if edStr := ed.PublicString(); edStr != expEd {
		t.Fatalf("Failed to generate EdDSA public string: %s", edStr)
	}

	if ecStr := ec.PublicString(); ecStr != expEc {
		t.Fatalf("Failed to generate ECDSA public string: %s", ecStr)
	}

	if rsStr := rs.PublicString(); rsStr != expRs {
		t.Fatalf("Failed to generate RSA public string: %s", rsStr)
	}
}

func TestEncode(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	expEd := "DHGQKw0oDDcMcZArDSgMNwxxkCsNKAw3DHGQKw0oDDe1+1s+xW4vzlPSPGN3OTEStdBKaW3SHjMRGJL5rk6IAA=="
	expEc := "MHcCAQEEIPb5VRJEX5ZNB1kxPOHKMrVgC1LQ39HgUlwzhgjGWtQRoAoGCCqGSM49AwEHoUQDQgAEbJD2pRVSiKjnrnkEY+Cs+0n3RcJjfhO5hzz8H3ilMgvshY0xdkBEuEwIpIJFYXUCladiye+hmWOgFxOen8EoTw=="
	expRs := "MIGsAgEAAiEAvvXLuX5q1GIZG1DaBD+UvrFKr+ip2gd9iAWt2PR5bHUCAwEAAQIhAIAeZTrkyEQKNRIQotVq2x4UX+WS2QnGukHhN4hbI/8BAhEA7XUgdS1JHhCoMflnQLyWlQIRAM3fKs4ZV2DxNVKPsnyxZmECEQDLCVxgdQFRAMMgP/XGh7plAhAZzXy59CcleVXrkSMXycxBAhEA2d7gYrSEfm4W8g99MQlmcg=="

	if edStr := ed.Encode(); edStr != expEd {
		t.Fatalf("Encoding EdDSA failed: %s", edStr)
	}

	if ecStr := ec.Encode(); ecStr != expEc {
		t.Fatalf("Encoding ECDSA failed: %s", ecStr)
	}

	if rsStr := rs.Encode(); rsStr != expRs {
		t.Fatalf("Encoding RSA failed: %s", rsStr)
	}
}

func TestDecodeHappy(t *testing.T) {
	ed := key.EdDSAKeyPair{}
	ec := key.ECDSAKeyPair{}
	rs := key.RSAKeyPair{}

	edStr := "DHGQKw0oDDcMcZArDSgMNwxxkCsNKAw3DHGQKw0oDDe1+1s+xW4vzlPSPGN3OTEStdBKaW3SHjMRGJL5rk6IAA=="
	ecStr := "MHcCAQEEIPb5VRJEX5ZNB1kxPOHKMrVgC1LQ39HgUlwzhgjGWtQRoAoGCCqGSM49AwEHoUQDQgAEbJD2pRVSiKjnrnkEY+Cs+0n3RcJjfhO5hzz8H3ilMgvshY0xdkBEuEwIpIJFYXUCladiye+hmWOgFxOen8EoTw=="
	rsStr := "MIGsAgEAAiEAvvXLuX5q1GIZG1DaBD+UvrFKr+ip2gd9iAWt2PR5bHUCAwEAAQIhAIAeZTrkyEQKNRIQotVq2x4UX+WS2QnGukHhN4hbI/8BAhEA7XUgdS1JHhCoMflnQLyWlQIRAM3fKs4ZV2DxNVKPsnyxZmECEQDLCVxgdQFRAMMgP/XGh7plAhAZzXy59CcleVXrkSMXycxBAhEA2d7gYrSEfm4W8g99MQlmcg=="

	if err := ed.Decode(edStr); err != nil || !ed.PrivateKey.Equal(buildEdDSAKey()) {
		t.Fatalf("Failed to decode EdDSAKeyPair from string: %v.", err)
	}

	if err := ec.Decode(ecStr); err != nil || !ec.PrivateKey.Equal(buildECDSAKey()) {
		t.Fatalf("Failed to decode ECDSAKeyPair from string: %v.", err)
	}

	if err := rs.Decode(rsStr); err != nil || !rs.PrivateKey.Equal(buildRSAKey()) {
		t.Fatalf("Failed to decode RSAKeyPair from string: %v.", err)
	}
}

func TestDecodeBadBase64Encoding(t *testing.T) {
	ed := key.EdDSAKeyPair{}
	ec := key.ECDSAKeyPair{}
	rs := key.RSAKeyPair{}

	badBase64 := "^^^$$$####^&&@"

	if err := ed.Decode(badBase64); err == nil {
		t.Fatal("Decoding bad base64 EdDSA string did not produce an error")
	}

	if err := ec.Decode(badBase64); err == nil {
		t.Fatal("Decoding bad base64 ECDSA string did not produce an error")
	}

	if err := rs.Decode(badBase64); err == nil {
		t.Fatal("Decoding bad base64 RSA string did not produce an error")
	}
}

func TestDecodeBadCertEncoding(t *testing.T) {
	ec := key.ECDSAKeyPair{}
	rs := key.RSAKeyPair{}

	badCert := "abcdefGHIJKabc=="

	if err := ec.Decode(badCert); err == nil {
		t.Fatal("Decoding bad ECDSA string did not produce an error")
	}

	if err := rs.Decode(badCert); err == nil {
		t.Fatal("Decoding bad RSA string did not produce an error")
	}

}

func TestKeyHappy(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	if !ed.Key().Equal(ed.PrivateKey) { t.Fatalf("Key returned by EdDSA Key() was not private key") }
	if !ec.Key().Equal(ec.PrivateKey) { t.Fatalf("Key returned by ECDSA Key() was not private key") }
	if !rs.Key().Equal(rs.PrivateKey) { t.Fatalf("Key returned by RSA Key() was not private key") }
}

func TestPublicKeyHappy(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	if !ed.PublicKey().(ed25519.PublicKey).Equal(ed.PrivateKey.Public()) { t.Fatalf("Key returned by EdDSA Key() was not private key") }
	if !ec.PublicKey().(*ecdsa.PublicKey).Equal(ec.PrivateKey.Public()) { t.Fatalf("Key returned by ECDSA Key() was not private key") }
	if !rs.PublicKey().(*rsa.PublicKey).Equal(rs.PrivateKey.Public()) { t.Fatalf("Key returned by RSA Key() was not private key") }
}

func TestSigningMethodHappy(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	if ed.SigningMethod() != jwt.SigningMethodEdDSA { t.Fatalf("EdDSA signing method was not jwt.SigningMethodEdDSA") }

	if ec.SigningMethod() != jwt.SigningMethodES256 { t.Fatalf("ECDSA with 256 bits signing method was not jwt.SigningMethodES256") }
	ec.NumBits = 512
	if ec.SigningMethod() != jwt.SigningMethodES512 { t.Fatalf("ECDSA with 512 bits signing method was not jwt.SigningMethodES512") }
	ec.NumBits = 384
	if ec.SigningMethod() != jwt.SigningMethodES384 { t.Fatalf("ECDSA with 384 bits signing method was not jwt.SigningMethodES384") }
	ec.NumBits = 1
	if ec.SigningMethod() != jwt.SigningMethodES256 { t.Fatalf("ECDSA with bad bits signing method did not default to jwt.SigningMethodES256") }

	if rs.SigningMethod() != jwt.SigningMethodPS256 { t.Fatalf("RSA with 256 bits signing method was not jwt.SigningMethodRS256") }
	rs.NumBits = 384
	if rs.SigningMethod() != jwt.SigningMethodPS384 { t.Fatalf("RSA with 256 bits signing method was not jwt.SigningMethodRS384") }
	rs.NumBits = 512
	if rs.SigningMethod() != jwt.SigningMethodPS512 { t.Fatalf("RSA with 512 bits signing method was not jwt.SigningMethodRS512") }
	rs.NumBits = 1
	if rs.SigningMethod() != jwt.SigningMethodPS256 { t.Fatalf("RSA with bad bits signing method did not default to jwt.SigningMethodRS256") }
}

func TestSetExpiresHappy(t *testing.T) {
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey() }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey() }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey() }

	later := time.Now().Add(time.Hour)

	ed.SetExpires(later)
	if ed.Expires != later { t.Fatalf("Failed to set expires on EdDSA") }

	ec.SetExpires(later)
	if ec.Expires != later { t.Fatalf("Failed to set expires on ECDSA") }

	rs.SetExpires(later)
	if rs.Expires != later { t.Fatalf("Failed to set expires on RSA") }
}

func TestExpiresAtHappy(t *testing.T) {
	later := time.Now().Add(time.Hour)
	ed := &key.EdDSAKeyPair{ PrivateKey: buildEdDSAKey(), KeyMeta: key.KeyMeta{ Expires: later} }
	ec := &key.ECDSAKeyPair{ NumBits: 256, PrivateKey: buildECDSAKey(), KeyMeta: key.KeyMeta{ Expires: later} }
	rs := &key.RSAKeyPair{ NumBits: 256, PrivateKey: buildRSAKey(), KeyMeta: key.KeyMeta{ Expires: later} }

	if ed.ExpiresAt() != later { t.Fatalf("Failed to get expire time on EdDSA") }
	if ec.ExpiresAt() != later { t.Fatalf("Failed to get expire time on ECDSA") }
	if rs.ExpiresAt() != later { t.Fatalf("Failed to get expire time on RSA") }
}
