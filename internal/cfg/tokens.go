package cfg

type Tokens struct {
	// One of RSA, ECDSA, or EdDSA
	Algorithm     string `json:"algorithm"`
	// Only applies for RSA and ECDSA
	NumBits       int    `json:"num_bits"`
	// Whether or not to save the private keys in the database
	PersistKeys   bool   `json:"persist_keys"`
	// How long to keep signing keys alive
	KeyDuration   string `json:"key_duration"`
	// How long to issue tokens for
	TokenDuration string `json:"token_duration"`
  // Claims to give all tokens
	Issuer   string   `json:"issuer"`
	Subject  string   `json:"subject"`
	Audience []string `json:"audience"`
}
