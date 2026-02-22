package key_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/base64"
	"stoke/internal/ent"
	"stoke/internal/key"
	"stoke/internal/testutil"
	"testing"
	"time"

	"hppr.dev/stoke"
)

var kc = key.PrivateKeyCache[ed25519.PrivateKey]{
	KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ edKeyPair },
	Ctx:           testutil.NewMockContext(),
	KeyDuration:   time.Hour,
	TokenDuration: time.Minute,
	PersistKeys:   false,
}


func TestPrivateKeyCachePublicKeys(t *testing.T) {
	publicJWKset, err := kc.PublicKeys(kc.Ctx)
	if err != nil {
		t.Fatalf("Error getting PrivateKeyCache public key json: %v", err)
	}

	expJWKStr := `{"exp":"0001-01-01T00:00:00.1Z","keys":[{"kty":"OKP","use":"sig","kid":"p-0","crv":"ed25519","x":"tftbPsVuL85T0jxjdzkxErXQSmlt0h4zERiS-a5OiAA="}]}`

	if string(publicJWKset) != expJWKStr {
		t.Logf("Public keys does not match expected value: \n%s\n%s", string(publicJWKset), expJWKStr)
		t.Fail()
	}
}

func TestPrivateKeyCacheGenerate(t *testing.T) {
	genCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ edKeyPair },
		Ctx:           testutil.NewMockContext(),
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

func TestPrivateKeyCacheGeneratePersistsKeys(t *testing.T) {
	ctx := testutil.NewMockContext(testutil.WithDatabase(t))
	genCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ edKeyPair },
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   true,
	}
	if err := genCache.Generate(ctx) ; err != nil {
		t.Fatalf("An error occured while generating a PrivateKey: %v", err)
	}
	if count, err := ent.FromContext(ctx).PrivateKey.Query().Count(ctx); err != nil || count != 1 {
		t.Fatalf("Did not persist generated keys when persist keys enabled: error %v, key count %d", err, count)
	}
}

func TestPrivateKeyCacheGenerateDoesNotReturnErrorIfPersistFails(t *testing.T) {
	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t, testutil.ReturnsMutateErrors()),
	)

	genCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ edKeyPair },
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   true,
	}
	if err := genCache.Generate(ctx) ; err != nil {
		t.Fatalf("An error occured while generating a PrivateKey: %v", err)
	}
	if len(genCache.KeyPairs) != 2 {
		t.Fatal("KeyPairs length was not 2 after generate")
	}
}

func TestPrivateKeyCacheGenerateWithEmptyKeyPairs(t *testing.T) {
	ctx := testutil.NewMockContext()
	genCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{},
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   false,
	}
	if err := genCache.Generate(ctx) ; err == nil {
		t.Fatal("Generated key when no keypairs existed")
	}
}

func TestPrivateKeyCacheGenerateFailsOnKeyPairGenerateFail(t *testing.T) {
	ctx := testutil.NewMockContext()
	genCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ BadKeyPair{} },
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   false,
	}
	if err := genCache.Generate(ctx) ; err == nil {
		t.Fatal("Generated key when current key pair was bad")
	}
}

func TestPrivateKeyCacheBootstrap(t *testing.T) {
	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t, testutil.ForeverKey()),
	)

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

	if len(bsCache.KeyPairs) != 1 {
		t.Logf("Do not have expected number of bootstrapped keys: %d", len(bsCache.KeyPairs))
		t.FailNow()
	}

	if !bsCache.KeyPairs[0].Key().Equal(edKey) {
		t.Logf("Bootstrapped key does not match:\nExp %s\nRes %s", base64.URLEncoding.EncodeToString(edKey), base64.URLEncoding.EncodeToString(bsCache.KeyPairs[0].Key()))
		t.Fail()
	}
}

func TestPrivateKeyCacheBootstrapCreatesNewKeyIfExpired(t *testing.T) {
	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t, testutil.KeyWithExpires(time.Now().Add(-time.Hour))),
	)

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

	if len(bsCache.KeyPairs) != 1 {
		t.Logf("Do not have expected number of bootstrapped keys: %d", len(bsCache.KeyPairs))
		t.FailNow()
	}

	if bsCache.KeyPairs[0].Key().Equal(edKey) {
		t.Logf("Bootstrapped key matches when it should have been regenerated:\nRes %s", base64.URLEncoding.EncodeToString(bsCache.KeyPairs[0].Key()))
		t.Fail()
	}
}

func TestPrivateKeyCacheBootstrapReturnsAnErrorIfPairFailsToGenerate(t *testing.T) {
	ctx := testutil.NewMockContext(testutil.WithDatabase(t))
	bsCache := key.PrivateKeyCache[ed25519.PrivateKey]{
		KeyPairs:      []key.KeyPair[ed25519.PrivateKey]{ &BadKeyPair{} },
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   true,
	}
	if err := bsCache.Bootstrap(ctx, &BadKeyPair{}); err == nil {
		t.Log("Was able to bootstrap cache with bad key text")
		t.Fail()
	}
}

func TestPrivateKeyCacheBootstrapReturnsAnErrorIfPairFailsToDecode(t *testing.T) {
	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t, testutil.ForeverKeyWithText("baddad1234==")),
	)

	bsCache := key.PrivateKeyCache[*ecdsa.PrivateKey]{
		KeyPairs:      []key.KeyPair[*ecdsa.PrivateKey]{},
		Ctx:           ctx,
		KeyDuration:   time.Hour,
		TokenDuration: time.Minute,
		PersistKeys:   true,
	}
	if err := bsCache.Bootstrap(ctx, &key.ECDSAKeyPair{}); err == nil {
		t.Log("Was able to bootstrap cache with bad key text")
		t.Fail()
	}
}

func TestPrivateKeyCacheClean(t *testing.T) {
	expiredTime := time.Now().Add(-time.Minute) 
	okTime := time.Now().Add(time.Minute) 

	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t,
			testutil.KeyWithExpires(expiredTime),
			testutil.KeyWithExpires(okTime),
		),
	)

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
		t.FailNow()
	}

	if cache.KeyPairs[0] != k2 {
		t.Log("Clean key pair is not expected value")
		t.Fail()
	}

	dbKeys := ent.FromContext(ctx).PrivateKey.Query().AllX(ctx)

	if len(dbKeys) != 1 {
		t.Logf("Number of keys in the database did not match expeced value: expected:2, actual: %d", len(dbKeys))
		t.FailNow()
	}

	if !dbKeys[0].Expires.Equal(okTime) {
		t.Logf("Expected key left in database to have unexpired time(%v): %v", okTime , dbKeys[0])
		t.Fail()
	}
}

func TestPrivateKeyCacheCleanLogsErrorsIfDatabaseFailsDelete(t *testing.T) {
	expiredTime := time.Now().Add(-time.Minute) 
	okTime := time.Now().Add(time.Minute) 

	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t,
			testutil.KeyWithExpires(expiredTime),
			testutil.KeyWithExpires(okTime),
			testutil.ReturnsMutateErrors(),
		),
	)

	k1 := &key.EdDSAKeyPair{ PrivateKey: edKey }
	k1.SetExpires(expiredTime)
	k2 := &key.EdDSAKeyPair{ PrivateKey: edKey }
	k2.SetExpires(okTime)

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
		t.FailNow()
	}

	if cache.KeyPairs[0] != k2 {
		t.Log("Clean key pair is not expected value")
		t.Fail()
	}
}

func TestPrivateKeyCacheParseClaims(t *testing.T) {
	ctx := testutil.NewMockContext(testutil.WithDatabase(t))

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

func TestPrivateKeyCacheParseClaimsWithInvalidToken(t *testing.T) {
	ctx := testutil.NewMockContext(testutil.WithDatabase(t))

	tokenStr := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZSIsInN1YiI6IbaddadadVsZiIsImF1ZCI6WyJleWUiXSwiZXhwIjo5NTYxODM5ODIxMCwianRpIjoiazEiLCJoZWxsbyI6IndvcmxkIiwiZm9vIjoiYmFyIn0.6ZTrIrOHhUIvT5-3h2WGCwW0DCnuAJMPNdNIG5VMPWPgEix4fTqTUK8qsJUZH1SXbv0xmztPZOvvfuuykR06DQ"

	_, err := kc.ParseClaims(ctx, tokenStr, stoke.RequireToken() )
	if err == nil {
		t.Logf("Parse claims validated bad token: %v", err)
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

func TestNewPrivateKeyCacheWithManagementHappy(t *testing.T) {
	tokenDuration := 10 * time.Millisecond
	keyDuration := 100 * time.Millisecond

	ctx := testutil.NewMockContext(
		testutil.WithDatabase(t, testutil.KeyWithExpires(time.Now().Add(keyDuration))),
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	keyCache, err := key.NewPrivateKeyCache(tokenDuration, keyDuration, false, edKeyPair, ctx, "")
	if err != nil {
		t.Logf("Failed to create private key cache: %v", err)
		t.Fail()
	}

	if len(keyCache.KeyPairs) != 1 {
		t.Logf("Private key cache did not have correct number of keys before starting: %d", len(kc.KeyPairs))
		t.Fail()
	}

	currKey := keyCache.CurrentKey()

	time.Sleep(keyDuration - (tokenDuration * 2))
	time.Sleep(2 * time.Millisecond)

	// Should have created a new key and not activated it yet
	if len(keyCache.KeyPairs) != 2 {
		t.Logf("Private key cache did not have correct number of keys after creating: %d", len(kc.KeyPairs))
		t.Fail()
	}

	if keyCache.CurrentKey() != currKey {
		t.Log("Current key was updated before intended")
		t.Fail()
	}

	time.Sleep(tokenDuration)

	// Should have activated the new key but not cleaned the old one
	if len(keyCache.KeyPairs) != 2 {
		t.Logf("Private key cache did not have correct number of keys after activation: %d", len(kc.KeyPairs))
		t.Fail()
	}

	if keyCache.CurrentKey() == currKey {
		t.Log("New key was not activated")
		t.Fail()
	}

	currKey = keyCache.CurrentKey()

	time.Sleep(tokenDuration)

	// Should have cleaned the inactive key
	if len(keyCache.KeyPairs) != 1 {
		t.Logf("Private key cache did not have correct number of keys after clean: %d", len(kc.KeyPairs))
		t.Fail()
	}

	if keyCache.CurrentKey() != currKey {
		t.Log("Activated key was changed after activation")
		t.Fail()
	}

}
