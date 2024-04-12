package key

import (
	"context"
	"encoding/base64"
	"stoke/client/stoke"
	"stoke/internal/tel"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

type Claims struct {
	jwt.RegisteredClaims
	StokeClaims map[string]string `json:"stk"`
}

type TokenIssuer interface {
	IssueToken(Claims, context.Context) (string, string, error)
	RefreshToken(*jwt.Token, string, time.Duration, context.Context) (string, string, error)
	PublicKeys(context.Context) ([]byte, error)
	stoke.PublicKeyStore
}

type AsymetricTokenIssuer[P PrivateKey]  struct {
	Ctx context.Context
	KeyCache[P]
}

func (a *AsymetricTokenIssuer[P]) IssueToken(claims Claims, ctx context.Context) (string, string, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "AsymetricTokenIssuer.IssueToken")
	defer span.End()

	curr := a.CurrentKey()
	token, err := jwt.NewWithClaims(curr.SigningMethod(), claims).SignedString(curr.Key())
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to sign auth token")
		return "", "", err
	}

	refresh, err := curr.SigningMethod().Sign(token, curr.Key())
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
	logger := zerolog.Ctx(ctx)
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

	now := time.Now()
	stokeClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(extendTime))

	return a.IssueToken(Claims {
		RegisteredClaims: stokeClaims.RegisteredClaims,
		StokeClaims : stokeClaims.StokeClaims,
	}, ctx)
}

func (a *AsymetricTokenIssuer[P]) verifyRefreshToken(jwtToken *jwt.Token, refreshBytes []byte) error {
	var err error
	for _, curr := range a.keys {
		err = curr.SigningMethod().Verify(jwtToken.Raw, refreshBytes, curr.PublicKey())
		if err == nil {
			return nil
		}
	}
	return err
}
