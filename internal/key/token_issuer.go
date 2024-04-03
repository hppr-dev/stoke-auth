package key

import (
	"encoding/base64"
	"stoke/client/stoke"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	StokeClaims map[string]string `json:"stk"`
}

type TokenIssuer interface {
	IssueToken(Claims) (string, string, error)
	RefreshToken(*jwt.Token, string, time.Duration) (string, string, error)
	PublicKeys() ([]byte, error)
	stoke.PublicKeyStore
}

type AsymetricTokenIssuer[P PrivateKey]  struct {
	KeyCache[P]
}

func (a *AsymetricTokenIssuer[P]) IssueToken(claims Claims) (string, string, error) {
	curr := a.CurrentKey()
	token, err := jwt.NewWithClaims(curr.SigningMethod(), claims).SignedString(curr.Key())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to sign auth token")
		return "", "", err
	}
	refresh, err := curr.SigningMethod().Sign(token, curr.Key())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to sign refresh token")
		return "", "", err
	}
	return token, base64.StdEncoding.EncodeToString(refresh), err
}

func (a *AsymetricTokenIssuer[P]) RefreshToken(jwtToken *jwt.Token, refreshToken string, extendTime time.Duration) (string, string, error) {
	refreshBytes, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		logger.Error().
			Err(err).
			Str("refreshToken", refreshToken).
			Msg("Failed to decode refresh token")
		return "", "", err
	}

	if err := a.verifyRefreshToken(jwtToken, refreshBytes); err != nil {
		logger.Debug().
			Err(err).
			Str("refreshToken", refreshToken).
			Str("authToken", jwtToken.Raw).
			Msg("Failed to verify refresh token")
		return "", "", err
	}

	stokeClaims, ok := jwtToken.Claims.(*stoke.Claims)
	if !ok {
		logger.Debug().
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
	})
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
