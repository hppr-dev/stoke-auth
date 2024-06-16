package stoke

import (
	"errors"
	"regexp"
	"strings"
)

// Creates a token to require claims
func RequireToken() *Claims {
	return &Claims{
		StokeClaims: make(map[string]string),
	}
}

// Validates custom claim requirements.
// Automatically called when parsing token
func (c Claims) Validate() error {
	// If any of the alternate claims successfully validate, this claim is valid as well
	for _, other := range c.alternateClaims {
		other.StokeClaims = c.StokeClaims
		if err := other.Validate(); err == nil {
			return nil
		}
	}

	validated := true
	for _, pred := range c.requiredClaimPreds {
		validated = pred(c) && validated
		if !validated {
			return errors.New("Claims do not conform to claim requirements")
		}
	}
	return nil
}

// Requires that a token has a claim key
func (c *Claims) WithClaim(shortName, value string) *Claims {
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		testClaim(shortName, func(claimValue string) bool {
			return claimValue == value
		}),
	)
	return c
}

// Requires that a token has a claim key that refers to a value that matches a regex
func (c *Claims) WithClaimMatch(shortName, regex string) *Claims {
	compiledRegex, _ := regexp.Compile(regex)
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		testClaim(shortName, func(claimValue string) bool {
			return compiledRegex.MatchString(claimValue)
		}),
	)
	return c
}

// Requires that a token has a claim key that refers to a value that contains a substring
func (c *Claims) WithClaimContains(shortName, substring string) *Claims {
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		testClaim(shortName, func(claimValue string) bool {
			return strings.Contains(claimValue, substring)
		}),
	)
	return c
}

// Requires the current token claims or another token claims to validate
func (c *Claims) Or(other *Claims) *Claims {
	c.alternateClaims = append(c.alternateClaims, other)
	return c
}

// check if a claim exists and returns whether any of the comma separated keys pass a test
func testClaim(key string, test func(string) bool) func(Claims) bool {
	return func( inner Claims ) bool {
		val, ok := inner.StokeClaims[key]
		if !ok {
			return false
		}
		return mapOr(val, test)
	}
}

// splits a string by "," and tests that at least one predicate is true
func mapOr(s string, test func(string) bool) bool {
	for _, p := range strings.Split(s, ",") {
		if test(p) {
			return true
		}
	}
	return false
}
