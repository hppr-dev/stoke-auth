package key

import (
	"context"
	"encoding/json"
	"net/http"
	"stoke/internal/cluster"
	"stoke/internal/tel"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hppr.dev/stoke"
)

// FederatedTokenIssuer wraps a TokenIssuer and exposes a merged JWKS from the inner
// issuer and all discovered peers. It implements key.TokenIssuer and stoke.PublicKeyStore
// so that PublicKeys returns merged keys and ParseClaims verifies tokens using the merged set.
type FederatedTokenIssuer struct {
	Inner      TokenIssuer
	Discoverer cluster.Discoverer
	HTTPClient *http.Client
	BasePath   string
	RefreshSec int

	mu         sync.RWMutex
	cacheBytes []byte
	cacheExpiry time.Time
}

// NewFederatedTokenIssuer returns a TokenIssuer that delegates IssueToken, RefreshToken,
// and WithContext to inner, and returns merged JWKS from inner plus peers for PublicKeys
// and verifies tokens with the merged set in ParseClaims.
func NewFederatedTokenIssuer(inner TokenIssuer, discoverer cluster.Discoverer, httpClient *http.Client, basePath string, refreshSec int) TokenIssuer {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &FederatedTokenIssuer{
		Inner:       inner,
		Discoverer:  discoverer,
		HTTPClient:  httpClient,
		BasePath:    basePath,
		RefreshSec:  refreshSec,
	}
}

// IssueToken delegates to Inner.
func (f *FederatedTokenIssuer) IssueToken(claims *stoke.Claims, ctx context.Context) (string, string, error) {
	return f.Inner.IssueToken(claims, ctx)
}

// RefreshToken delegates to Inner.
func (f *FederatedTokenIssuer) RefreshToken(jwtToken *jwt.Token, refreshToken string, extendTime time.Duration, ctx context.Context) (string, string, error) {
	return f.Inner.RefreshToken(jwtToken, refreshToken, extendTime, ctx)
}

// WithContext stores this federated issuer in context so the web layer (IssuerFromCtx as PublicKeyStore)
// uses the merged set for PublicKeys and ParseClaims.
func (f *FederatedTokenIssuer) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, issuerCtxKey{}, f)
}

// PublicKeys returns the merged JWKS from Inner and all peers. Optionally cached for RefreshSec.
func (f *FederatedTokenIssuer) PublicKeys(ctx context.Context) ([]byte, error) {
	_, span := tel.GetTracer().Start(ctx, "FederatedTokenIssuer.PublicKeys")
	defer span.End()

	if f.RefreshSec > 0 {
		f.mu.RLock()
		if len(f.cacheBytes) > 0 && time.Now().Before(f.cacheExpiry) {
			b := make([]byte, len(f.cacheBytes))
			copy(b, f.cacheBytes)
			f.mu.RUnlock()
			return b, nil
		}
		f.mu.RUnlock()
	}

	merged, err := f.getMergedJWKSBytes(ctx)
	if err != nil {
		return nil, err
	}

	if f.RefreshSec > 0 {
		f.mu.Lock()
		f.cacheBytes = merged
		f.cacheExpiry = time.Now().Add(time.Duration(f.RefreshSec) * time.Second)
		f.mu.Unlock()
	}

	return merged, nil
}

// ParseClaims verifies the token using the merged key set so that tokens issued by other replicas validate.
func (f *FederatedTokenIssuer) ParseClaims(ctx context.Context, token string, claims *stoke.Claims, parserOpts ...jwt.ParserOption) (*jwt.Token, error) {
	_, span := tel.GetTracer().Start(ctx, "FederatedTokenIssuer.ParseClaims")
	defer span.End()

	merged, err := f.getMergedJWKSBytes(ctx)
	if err != nil {
		return nil, err
	}

	var jwks stoke.JWKSet
	if err := json.Unmarshal(merged, &jwks); err != nil {
		return nil, err
	}

	pkeys := make([]jwt.VerificationKey, 0, len(jwks.Keys))
	for _, k := range jwks.Keys {
		if k == nil {
			continue
		}
		pub, err := k.ToPublicKey()
		if err != nil {
			continue
		}
		pkeys = append(pkeys, pub)
	}

	keySet := jwt.VerificationKeySet{Keys: pkeys}
	keyfunc := func(*jwt.Token) (interface{}, error) {
		return keySet, nil
	}
	return jwt.ParseWithClaims(token, claims.New(), keyfunc, parserOpts...)
}

// getMergedJWKSBytes returns merged JWKS from Inner and peers. Caller may hold no locks.
func (f *FederatedTokenIssuer) getMergedJWKSBytes(ctx context.Context) ([]byte, error) {
	localJWKS, err := f.Inner.PublicKeys(ctx)
	if err != nil {
		return nil, err
	}
	peerURLs, err := f.Discoverer.Peers(ctx)
	if err != nil {
		return nil, err
	}
	return cluster.MergeJWKS(localJWKS, peerURLs, f.HTTPClient)
}
