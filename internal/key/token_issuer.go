package key

import (
	"context"
	"encoding/base64"
	"fmt"
	"stoke/internal/tel"
	"time"

	"hppr.dev/stoke"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

func IssuerFromCtx(ctx context.Context) TokenIssuer {
	return ctx.Value("issuer").(TokenIssuer)
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
	return context.WithValue(ctx, "issuer", a)
}

func (a *AsymetricTokenIssuer[P]) IssueToken(claims *stoke.Claims, ctx context.Context) (string, string, error) {
	logger := zerolog.Ctx(ctx).With().Str("function", "AsymetricTokenIssuer.IssueToken").Logger()
	ctx, span := tel.GetTracer().Start(ctx, "AsymetricTokenIssuer.IssueToken")
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
	priv := curr.Key()
	a.ReadUnlock()
	token, err := jwt.NewWithClaims(curr.SigningMethod(), claims).SignedString(priv)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to sign auth token")
		return "", "", err
	}

	refresh, err := curr.SigningMethod().Sign(token, priv)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to sign refresh token")
		return "", "", err
	}
	return token, base64.StdEncoding.EncodeToString(refresh), err
}

func (a *AsymetricTokenIssuer[P]) RefreshToken(jwtToken *jwt.Token, refreshToken string, extendTime time.Duration, ctx context.Context) (string, string, error) {
	logger := zerolog.Ctx(ctx).With().Str("function", "AsymetricTokenIssuer.RefreshToken").Logger()
	ctx, span := tel.GetTracer().Start(ctx, "AsymetricTokenIssuer.RefreshToken")
	defer span.End()

	refreshBytes, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("refreshToken", refreshToken).
			Msg("Failed to decode refresh token")
		return "", "", err
	}

	if err := a.verifyRefreshToken(jwtToken, refreshBytes); err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("refreshToken", refreshToken).
			Str("authToken", jwtToken.Raw).
			Msg("Failed to verify refresh token")
		return "", "", err
	}

	stokeClaims, ok := jwtToken.Claims.(*stoke.Claims)
	if !ok {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("refreshToken", refreshToken).
			Str("authToken", jwtToken.Raw).
			Type("claimsType", jwtToken.Claims).
			Interface("claimsValues", jwtToken.Claims).
			Msg("Failed to convert jwt.Claims to stoke.Claims")
		return "", "", err
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
			oldJwtID, _ = claims.StokeClaims[a.TokenRefreshCountKey]
		}
		
		if oldJwtID == "" {
			jwtID = fmt.Sprintf(jwtFormat, 1)
		} else {
			var gen int
			fmt.Sscanf(claims.ID, jwtFormat, &gen)
			if gen > a.TokenRefreshLimit {
				return fmt.Errorf("Token refresh limit reached.")
			}
			jwtID = fmt.Sprintf(jwtFormat, gen+1)
		}

		if a.TokenRefreshCountKey == "" {
			claims.ID = jwtID
		} else {
			claims.StokeClaims[a.TokenRefreshCountKey] = jwtID
		}
	}

	return nil
}
