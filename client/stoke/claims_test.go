package stoke_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

func TestClaimsJSONMarshalling(t *testing.T) {
	now := time.Now()
	expires := now.Add(time.Hour)
	notBefore := now.Add(-time.Hour)
	issued := now.Add(-time.Hour)

	c := &stoke.Claims{
		StokeClaims:      map[string]string{
			"": "should omit",
			"onefish" : "twofish",
			"3fish": "bluefish",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   "subject",
			Audience:  []string{"aud1", "aud2"},
			ExpiresAt: jwt.NewNumericDate(expires),
			NotBefore: jwt.NewNumericDate(notBefore),
			IssuedAt:  jwt.NewNumericDate(issued),
			ID:        "12345asdf",
		},
	}

	marshalled, err := json.Marshal(c)
	if err != nil {
		t.Logf("An error occurred while marshalling: %v", err)
		t.FailNow()
	}

	toMap := make(map[string]interface{})
	if err = json.Unmarshal(marshalled, &toMap); err != nil {
		t.Logf("An error occurred while unmarshalling: %v", err)
		t.FailNow()
	}

	for k, v := range toMap {
		switch k {
		case "3fish":
			if val, ok := v.(string); !ok || val != "bluefish" {
				t.Logf("3fish did not match expected value. Expected bluefish, but got %v", v)
				t.Fail()
			}
		case "onefish":
			if val, ok := v.(string); !ok || val != "twofish" {
				t.Logf("onefish did not match expected value. Expected twofish, but got %v", v)
				t.Fail()
			}
		case "iss":
			if val, ok := v.(string); !ok || val != "issuer" {
				t.Logf("issuer did not match expected value. Expected issuer, but got %v", v)
				t.Fail()
			}
		case "sub":
			if val, ok := v.(string); !ok || val != "subject" {
				t.Logf("subject did not match expected value. Expected subject, but got %v", v)
				t.Fail()
			}
		case "aud":
			if val, ok := v.([]interface{}); !ok || len(val) != 2 || val[0] != "aud1" || val[1] != "aud2" {
				t.Logf("aud did not match expected value. Expected [aud1, aud2] , but got %T %v", v, v)
				t.Fail()
			}
		case "exp":
			if val, ok := v.(float64); !ok || val != float64(expires.Unix()) {
				t.Logf("exp did not match expected value. Expected %T %v, but got %T %v", expires, expires, v, v)
				t.Fail()
			}
		case "nbf":
			if val, ok := v.(float64); !ok || val != float64(notBefore.Unix()) {
				t.Logf("nbf did not match expected value. Expected %T %v, but got %T %v", notBefore, notBefore, v, v)
				t.Fail()
			}
		case "iat":
			if val, ok := v.(float64); !ok || val != float64(issued.Unix()) {
				t.Logf("iat did not match expected value. Expected %T %v, but got %T %v", issued, issued, v, v)
				t.Fail()
			}
		case "jti":
			if val, ok := v.(string); !ok || val != "12345asdf" {
				t.Logf("jti did not match expected value. Expected 12345asdf, but got %v", v)
				t.Fail()
			}
		default:
			t.Logf("Received unexpected key: %s : %T %v", k, v, v)
			t.Fail()
		}
	}
}

func TestClaimsJSONUnmarshalling(t *testing.T) {
	now := time.Now()
	expires := now.Add(time.Hour)
	notBefore := now.Add(-time.Hour)
	issued := now.Add(-time.Hour)
	claimJson := fmt.Sprintf(`{"jti":"1q2w3e4r","exp":%d,"nbf":%d,"iat":%d,"iss":"issuer","sub":"subject","aud":["foo","bar"],"hello":"world","foo":"bar"}`, expires.Unix(), notBefore.Unix(), issued.Unix())

	claims := &stoke.Claims{}
	if err := json.Unmarshal([]byte(claimJson), claims); err != nil {
		t.Logf("Unmarshalling produced an error: %v\n%s", err, claimJson)
		t.FailNow()
	}

	if ex, err := claims.GetExpirationTime(); err != nil || ex.Unix() != expires.Unix() {
		t.Logf("exp did not match: Expected %v, but got %v", expires, ex)
		t.Fail()
	}

	if is, err := claims.GetIssuedAt(); err != nil || is.Unix() != issued.Unix() {
		t.Logf("iat did not match: Expected %v, but got %v", issued, is)
		t.Fail()
	}

	if nb, err := claims.GetNotBefore(); err != nil || nb.Unix() != notBefore.Unix() {
		t.Logf("nbt did not match: Expected %v, but got %v", notBefore, nb)
		t.Fail()
	}

	if aud, err := claims.GetAudience(); err != nil || len(aud) != 2  || aud[0] != "foo" || aud[1] != "bar" {
		t.Logf("aud did not match: Expected [foo,bar], but got %v", aud)
		t.Fail()
	}

	if sub, err := claims.GetSubject(); err != nil || sub != "subject" {
		t.Logf("sub did not match: Expected subject, but got %v", sub)
		t.Fail()
	}

	if iss, err := claims.GetIssuer(); err != nil || iss != "issuer" {
		t.Logf("iss did not match: Expected subject, but got %v", iss)
		t.Fail()
	}
	
	if claims.ID != "1q2w3e4r" {
		t.Logf("jti did not match: Expected 1q2w3e4r, but got %s", claims.ID)
		t.Fail()
	}

	if v, ok := claims.StokeClaims["hello"]; !ok || v != "world" {
		t.Logf("hello did not match: Expected world, got %s", v)
		t.Fail()
	}

	if v, ok := claims.StokeClaims["foo"]; !ok || v != "bar" {
		t.Logf("hello did not match: Expected bar, got %s", v)
		t.Fail()
	}
}

func TestClaimsNewFromClaim(t *testing.T) {
	now := time.Now()
	expires := now.Add(time.Hour)
	notBefore := now.Add(-time.Hour)
	issued := now.Add(-time.Hour)

	c := &stoke.Claims{
		StokeClaims:      map[string]string{
			"onefish" : "twofish",
			"3fish": "bluefish",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   "subject",
			Audience:  []string{"aud1", "aud2"},
			ExpiresAt: jwt.NewNumericDate(expires),
			NotBefore: jwt.NewNumericDate(notBefore),
			IssuedAt:  jwt.NewNumericDate(issued),
			ID:        "12345asdf",
		},
	}

	newClaims := c.New()
	for k, v := range newClaims.StokeClaims {
		t.Logf("Non-empty stoke claims on new claim: %s, %v", k, v)
		t.Fail()
	}
}

func TestClaimsWithClaim(t *testing.T) {
	c := stoke.RequireToken().WithClaim("exact", "match")
	c.StokeClaims["exact"] = "match"

	if err := c.Validate(); err != nil {
		t.Logf("Matching claim produced validation error: %v", err)
		t.Fail()
	}

	c.StokeClaims["exact"] = "not a match"

	if err := c.Validate(); err == nil {
		t.Logf("Non matching claim did not produce validation error")
		t.Fail()
	}
}

func TestClaimsWithClaimMatch(t *testing.T) {
	c := stoke.RequireToken().WithClaimMatch("regex", "[0-9][abc]")
	c.StokeClaims["regex"] = "8b"

	if err := c.Validate(); err != nil {
		t.Logf("Matching claim produced validation error: %v", err)
		t.Fail()
	}

	c.StokeClaims["regex"] = "nope"

	if err := c.Validate(); err == nil {
		t.Logf("Non matching claim did not produce validation error")
		t.Fail()
	}
}

func TestClaimsWithClaimContains(t *testing.T) {
	c := stoke.RequireToken().WithClaimContains("contains", "in")
	c.StokeClaims["contains"] = "thing"

	if err := c.Validate(); err != nil {
		t.Logf("Matching claim produced validation error: %v", err)
		t.Fail()
	}

	c.StokeClaims["contains"] = "nope"

	if err := c.Validate(); err == nil {
		t.Logf("Non matching claim did not produce validation error")
		t.Fail()
	}
}

func TestClaimsWithClaimListPart(t *testing.T) {
	c := stoke.RequireToken().WithClaim("listed", "one")

	if err := c.Validate(); err == nil {
		t.Logf("Missing claim did not produced validation error")
		t.Fail()
	}

	c.StokeClaims["listed"] = "four,one,five"

	if err := c.Validate(); err != nil {
		t.Logf("Matching claim produced validation error: %v", err)
		t.Fail()
	}

	c.StokeClaims["listed"] = "nope,nada,five"

	if err := c.Validate(); err == nil {
		t.Logf("Non matching claim did not produce validation error")
		t.Fail()
	}
}

func TestClaimsWithClaimOr(t *testing.T) {
	c := stoke.RequireToken().WithClaim("role", "capt").Or(
		stoke.RequireToken().WithClaim("role", "crew"),
	)
	c.StokeClaims["role"] = "capt"

	if err := c.Validate(); err != nil {
		t.Logf("Matching claim (capt) produced validation error: %v", err)
		t.Fail()
	}

	c.StokeClaims["role"] = "crew"

	if err := c.Validate(); err != nil {
		t.Logf("Matching claim (crew) produced validation error: %v", err)
		t.Fail()
	}

	c.StokeClaims["role"] = "staff"

	if err := c.Validate(); err == nil {
		t.Logf("Non matching claim did not produce validation error")
		t.Fail()
	}
}
