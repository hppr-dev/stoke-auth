package key_test

import (
	"crypto/ed25519"
	"stoke/internal/key"
	"testing"
	"time"

	"hppr.dev/stoke"
)

var kc = key.PrivateKeyCache[ed25519.PrivateKey]{
	KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ edKeyPair },
	Ctx:           NewMockContext(),
	KeyDuration:   time.Hour,
	TokenDuration: time.Minute,
	PersistKeys:   false,
}


func TestPrivateKeyCachePublicKeys(t *testing.T) {
	publicJWKset, err := kc.PublicKeys(kc.Ctx)
	if err != nil {
		t.Fatalf("Error getting PrivateKeyCache public key json: %v", err)
	}

	expJWKStr := `{"exp":"0001-01-01T00:00:00.1Z","keys":[{"kty":"EC","use":"sig","kid":"p-0","crv":"ed25519","x":"tftbPsVuL85T0jxjdzkxErXQSmlt0h4zERiS-a5OiAA="}]}`

	if string(publicJWKset) != expJWKStr {
		t.Logf("Public keys does not match expected value: \n%s\n%s", string(publicJWKset), expJWKStr)
		t.Fail()
	}
}

func TestPrivateKeyCacheGenerate(t *testing.T) {
	genCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ edKeyPair },
		Ctx:           NewMockContext(),
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   false,
	}
	if err := genCache.Generate(kc.Ctx) ; err != nil {
		t.Fatalf("An error occured while generating a PrivateKey: %v", err)
	}
	if len(genCache.KeyPairs) != 2 {
		t.Fatal("KeyPairs length was not 2 after generate")
	}
}

func TestPrivateKeyCacheBootstrap(t *testing.T) {
	ctx := WithDatabase(t, NewMockContext())
	bsCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{},
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   true,
	}
	if err := bsCache.Bootstrap(ctx, &key.EdDSAKeyPair{}); err != nil {
		t.Logf("Failed to boostrap cache: %v", err)
		t.Fail()
	}

	if len(bsCache.KeyPairs) != 1 || bsCache.KeyPairs[0].Key().Equal(buildEdDSAKey()) {
		t.Log("Bootstrapped keys do not match expected values")
		t.Fail()
	}
}

func TestPrivateKeyCacheClean(t *testing.T) {
	ctx := WithDatabase(t, NewMockContext())
	k1 := &key.EdDSAKeyPair{ PrivateKey: edKey }
	k1.SetExpires(time.Now().Add(-time.Minute))
	k2 := &key.EdDSAKeyPair{ PrivateKey: edKey }
	k2.SetExpires(time.Now().Add(time.Minute))

	cache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ k1, k2 },
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   true,
	}

	cache.Clean(ctx)

	if len(cache.KeyPairs) != 1 {
		t.Log("Clean did not remove expired certificates")
		t.Fail()
	}

	if cache.KeyPairs[0] != k2 {
		t.Log("Clean key pair is not expected value")
		t.Fail()
	}
}

func TestPrivateKeyCacheParseClaims(t *testing.T) {
	ctx := WithDatabase(t, NewMockContext())

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6Ik15c2VsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"

	jwtToken, err := kc.ParseClaims(ctx, tokenStr, stoke.RequireToken() )
	if err != nil {
		t.Logf("Parse claims returned an error: %v", err)
		t.Fail()
	}

	stokeClaims, ok := jwtToken.Claims.(*stoke.Claims)
	if !ok {
		t.Log("Could not convert claims to stoke.Claims")
		t.Fail()
	}

	if stokeClaims.StokeClaims["hello"] != "world" || stokeClaims.StokeClaims["foo"] != "bar" {
		t.Logf("Resulting claims did not have expected keys: %v", stokeClaims.StokeClaims)
		t.Fail()
	}
}

func TestPrivateKeyCacheKeys(t *testing.T) {
	keys := kc.Keys()
	if len(keys) != 1 || !keys[0].Key().Equal(buildEdDSAKey()) {
		t.Fatalf("Key returned by keys does not match expected key: \n%#v\n%#v", keys[0].Key(), buildEdDSAKey() )
	}
}

func TestPrivateKeyCacheReadLock(t *testing.T) {
	// Make sure ReadLock doesnt cause an error
	kc.ReadLock()
	kc.ReadLock()
	kc.ReadLock()
	kc.ReadUnlock()
	kc.ReadUnlock()
	kc.ReadUnlock()
}

func TestPrivateKeyCacheCurrentKey(t *testing.T) {
	if !edKeyPair.Key().Equal(kc.CurrentKey().Key()) {
		t.Fatal("Returned KeyCache Key did not match CurrentKey")
	}
}

func TestPrivateKeyCacheInitManagement(t *testing.T) {

}
