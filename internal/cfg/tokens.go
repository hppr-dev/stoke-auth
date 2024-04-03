package cfg

import (
	"regexp"
	"strconv"
	"time"
)

type Tokens struct {
	// One of RSA, ECDSA, or EdDSA
	Algorithm        string `json:"algorithm"`
	// Only applies for RSA and ECDSA
	NumBits          int    `json:"num_bits"`
	// Whether or not to save the private keys in the database
	PersistKeys      bool   `json:"persist_keys"`
	// How long to keep signing keys alive
	KeyDurationStr   string `json:"key_duration"`
	// How long to issue tokens for
	TokenDurationStr string `json:"token_duration"`
	// Include Not Before Header
	IncludeNotBefore bool     `json:"include_not_before"`
	// Include Issued At Header
	IncludeIssuedAt bool     ` json:"include_issued_at"`
	// Issuer to set on all tokens
	Issuer           string   `json:"issuer"`
	// Subject to set on all tokens
	Subject          string   `json:"subject"`
	// Audience to set on all tokens
	Audience         []string `json:"audience"`
	// List of user identifiers to include in Tokens
	// May have any or all of the following keys: username, first_name, last_name, full_name, email
	UserInfo         map[string]string `json:"user_info"`

	// Non-parsed fields
	TokenDuration time.Duration `json:"-"`
	KeyDuration time.Duration   `json:"-"`
}

func (t *Tokens) ParseDurations() {
	t.TokenDuration = strToDuration(t.TokenDurationStr)
	t.KeyDuration = strToDuration(t.KeyDurationStr)
}

var durationRegex *regexp.Regexp = regexp.MustCompile(`(\d+)([sSmMhHdDyY])`)

func strToDuration(s string) time.Duration {
	matches := durationRegex.FindStringSubmatch(s)
	if len(matches) != 3 {
		panic("Duration string did not match regex [0-9]+[sSmMhHdDyY")
	}
	num, _ := strconv.Atoi(matches[1])
	dur := time.Duration(num)
	switch matches[2] {
	case "s", "S":
		return time.Second * dur
	case "m", "M":
		return time.Minute * dur
	case "h", "H":
		return time.Hour * dur
	case "d", "D":
		return time.Hour * 24 * dur
	case "y", "Y":
		return time.Hour * 24 * 265 * dur
	}
	panic("Unreachable. If it reaches here, it means that durationRegex is broken.")
}
