package stoke

import (
	"errors"
	"regexp"
	"strings"
)

// Creates a token to require claims
func RequireToken() *Claims {
	return &Claims{}
}

// Validates custom claim requirements.
// Automatically called when parsing token
func (c Claims) Validate() error {
	// If any of the alternate claims successfully validate, this claim is valid as well
	for _, other := range c.alternateClaims {
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

// Requires that a token has a claim key that refers to a value that is equal to value
func (c *Claims) WithClaim(key, value string) *Claims {
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		func(inner Claims) bool {
			val, ok := inner.StokeClaims[key]
			return ok && val == value
		})
	return c
}

// Requires that a token has a claim key that refers to a value that matches a regex
func (c *Claims) WithClaimMatch(key, regex string) *Claims {
	compiledRegex, _ := regexp.Compile(regex)
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		func(inner Claims) bool {
			val, ok := inner.StokeClaims[key]
			return ok && compiledRegex.MatchString(val)
		})
	return c
}

// Requires that a token has a claim key that refers to a value that contains a substring
func (c *Claims) WithClaimContains(key, substring string) *Claims {
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		func(inner Claims) bool {
			val, ok := inner.StokeClaims[key]
			return ok && strings.Contains(val, substring)
		})
	return c
}

// Requires that a token has a claim key that refers to a comma seperated list with an item
func (c *Claims) WithClaimListPart(key, item string) *Claims {
	c.requiredClaimPreds = append(c.requiredClaimPreds,
		func(inner Claims) bool {
			val, ok := inner.StokeClaims[key]
			if !ok {
				return false
			}
			for _, p := range strings.Split(val, ",") {
				if p == item {
					return true
				}
			}
			return false
		})
	return c
}

// Requires the current token claims or another token claims to validate
func (c *Claims) Or(other Claims) *Claims {
	c.alternateClaims = append(c.alternateClaims, other)
	return c
}
