package key

import (
	"context"
	"encoding/base64"
	"fmt"
	"stoke/internal/tel"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"hppr.dev/stoke"
)

type issuerCtxKey struct{}

func IssuerFromCtx(ctx context.Context) TokenIssuer {
	return ctx.Value(issuerCtxKey{}).(TokenIssuer)
}

// localKeysOnlyCtxKey is the context key for the "local keys only" flag (e.g. from ?local=true on /api/pkeys).
type localKeysOnlyCtxKey struct{}

// LocalKeysOnly returns true if the context requests only this node's keys (no peer merge).
// Used to avoid recursion when a peer fetches our keys for JWKS merge.
func LocalKeysOnly(ctx context.Context) bool {
	v, _ := ctx.Value(localKeysOnlyCtxKey{}).(bool)
	return v
}

// WithLocalKeysOnly marks the context so PublicKeys returns only local keys (no peer merge).
func WithLocalKeysOnly(ctx context.Context) context.Context {
	return context.WithValue(ctx, localKeysOnlyCtxKey{}, true)
}

type TokenIssuer interface {
	IssueToken(*stoke.Claims, context.Context) (string, string, error)
	RefreshToken(*jwt.Token, string, time.Duration, context.Context) (string, string, error)
	PublicKeys(context.Context) ([]byte, error)
	WithContext(context.Context) context.Context
	stoke.PublicKeyStore
}

type AsymetricTokenIssuer[P PrivateKey]  struct {
	Ctx context.Context
	TokenRefreshLimit int
	TokenRefreshCountKey string

	KeyCache[P]
}

func (a *AsymetricTokenIssuer[P]) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, issuerCtxKey{}, a)
}

func (a *AsymetricTokenIssuer[P]) IssueToken(claims *stoke.Claims, ctx context.Context) (string, string, error) {
	logger := zerolog.Ctx(ctx).With().Str("function", "AsymetricTokenIssuer.IssueToken").Logger()
	_, span := tel.GetTracer().Start(ctx, "AsymetricTokenIssuer.IssueToken")
	defer span.End()

	if err := a.setJWTID(claims); err != nil {
		logger.Debug().
			Err(err).
			Int("refreshLimit", a.TokenRefreshLimit).
			Msg("Reached token refresh limit")
		return "", "", err
	}

	a.ReadLock()
	curr := a.CurrentKey()
	keyId := a.CurrentKeyId()
	priv := curr.Key()
	a.ReadUnlock()

	claims.StokeClaims["kid"] = keyId

	token := jwt.NewWithClaims(curr.SigningMethod(), claims)
	token.Header["kid"] = keyId
	signed, tok_err := token.SignedString(priv)
	if tok_err != nil {
		return "", "", tok_err
	}
	refresh, ref_err := curr.SigningMethod().Sign(signed, priv)
	return signed, base64.URLEncoding.EncodeToString(refresh), ref_err
}

func (a *AsymetricTokenIssuer[P]) RefreshToken(jwtToken *jwt.Token, refreshToken string, extendTime time.Duration, ctx context.Context) (string, string, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("function", "AsymetricTokenIssuer.RefreshToken").
		Str("refreshToken", refreshToken).
		Str("authToken", jwtToken.Raw).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "AsymetricTokenIssuer.RefreshToken")
	defer span.End()

	refreshBytes, err := base64.URLEncoding.DecodeString(refreshToken)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to decode refresh token")
		return "", "", err
	}

	if err := a.verifyRefreshToken(jwtToken, refreshBytes); err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to verify refresh token")
		return "", "", err
	}

	stokeClaims, ok := jwtToken.Claims.(*stoke.Claims)
	if !ok {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Type("claimsType", jwtToken.Claims).
			Interface("claimsValues", jwtToken.Claims).
			Msg("Failed to convert jwt.Claims to stoke.Claims")
		return "", "", fmt.Errorf("Failed to convert jwt.Claims")
	}

	oldTime := stokeClaims.RegisteredClaims.ExpiresAt
	stokeClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(oldTime.Add(extendTime))

	return a.IssueToken(stokeClaims, ctx)
}

func (a *AsymetricTokenIssuer[P]) verifyRefreshToken(jwtToken *jwt.Token, refreshBytes []byte) error {
	var err error
	a.ReadLock()
	defer a.ReadUnlock()
	for _, curr := range a.Keys() {
		err = curr.SigningMethod().Verify(jwtToken.Raw, refreshBytes, curr.PublicKey())
		if err == nil {
			return nil
		}
	}
	return err
}

const jwtFormat = "k%d"

func (a *AsymetricTokenIssuer[P]) setJWTID(claims *stoke.Claims) error {
	if a.TokenRefreshLimit != 0 {
		var oldJwtID, jwtID string
		if a.TokenRefreshCountKey == "" {
			oldJwtID = claims.ID
		} else {
			oldJwtID = claims.StokeClaims[a.TokenRefreshCountKey]
		}
		
		if oldJwtID == "" {
			jwtID = fmt.Sprintf(jwtFormat, a.TokenRefreshLimit)
		} else {
			var gen int
			_, err := fmt.Sscanf(oldJwtID, jwtFormat, &gen)
			if err != nil {
				return err
			}
			if gen == 0 {
				return fmt.Errorf("Token refresh limit reached.")
			}
			jwtID = fmt.Sprintf(jwtFormat, gen-1)
		}

		if a.TokenRefreshCountKey == "" {
			claims.ID = jwtID
		} else {
			claims.StokeClaims[a.TokenRefreshCountKey] = jwtID
		}
	}

	return nil
}
