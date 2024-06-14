package web

import (
	"context"
	"fmt"
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/ent/ogent"
	"stoke/internal/key"
	"stoke/internal/tel"
	"stoke/internal/usr"
	"time"

	"hppr.dev/stoke"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

type LoginApiHandler struct {}

// Login implements ogent.Handler.
func (h *entityHandler) Login(ctx context.Context, req *ogent.LoginReq) (ogent.LoginRes, error) {
	logger := zerolog.Ctx(ctx)

	ctx, span := tel.GetTracer().Start(ctx, "LoginHandler")
	defer span.End()

	user, claims, err := usr.ProviderFromCtx(ctx).GetUserClaims(req.Username, req.Password, true, ctx)
	if err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to get claims from provider")
		return &ogent.LoginUnauthorized{}, nil
	}
	
	claimMap := make(map[string]string)
	reqMet := make(map[string]bool)
	reqClaims := make(map[string]string)

	for _, claimReq := range req.RequiredClaims {
		reqMet[claimReq.Name] = false
		reqClaims[claimReq.Name] = claimReq.Value
	}

	for _, claim := range claims {
		if reqValue, exists := reqClaims[claim.ShortName]; exists && !reqMet[claim.ShortName] {
			reqMet[claim.ShortName] = reqValue == claim.Value
		}

		if value, exists := claimMap[claim.ShortName]; exists {
			claimMap[claim.ShortName] = value + "," + claim.Value
		} else {
			claimMap[claim.ShortName] = claim.Value
		}
	}

	for reqName, valueMatched := range reqMet {
		if !valueMatched {
			logger.Debug().
				Str("claimShortName", reqName).
				Str("requiredValue", reqClaims[reqName]).
				Msg("User did not have required claims.")
			return &ogent.LoginUnauthorized{}, nil
		}
	}

	populateUserInfo(cfg.Ctx(ctx), user, claimMap)

	token, refresh, err := key.IssuerFromCtx(ctx).IssueToken(&stoke.Claims{
		StokeClaims : claimMap,
		RegisteredClaims: createRegisteredClaims(cfg.Ctx(ctx).Tokens),
	}, ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not issue token")
		return &ogent.LoginBadRequest{}, nil
	}

	return &ogent.LoginOK{
		Token: token,
		Refresh: refresh,
	}, nil
}

func createRegisteredClaims(c cfg.Tokens) jwt.RegisteredClaims {
	now := time.Now()
	minClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(c.TokenDuration)),

		// Below fields are omitted if they are not included in config
		Issuer:    c.Issuer,
		Subject:   c.Subject,
		Audience:  c.Audience,
	}

	if c.IncludeNotBefore {
		minClaims.NotBefore = jwt.NewNumericDate(now)
	}

	if c.IncludeIssuedAt {
		minClaims.IssuedAt = jwt.NewNumericDate(now)
	}

	return minClaims
}

func populateUserInfo(c *cfg.Config, user *ent.User, t map[string]string) {
	usernameKey, ok := c.Tokens.UserInfo["username"]
	if ok {
		t[usernameKey] = user.Username
	}

	fnameKey, ok := c.Tokens.UserInfo["first_name"]
	if ok {
		t[fnameKey] = user.Fname
	}

	lnameKey, ok := c.Tokens.UserInfo["last_name"]
	if ok {
		t[lnameKey] = user.Lname
	}

	nameKey, ok := c.Tokens.UserInfo["full_name"]
	if ok {
		t[nameKey] = fmt.Sprintf("%s %s", user.Fname, user.Lname)
	}

	emailKey, ok := c.Tokens.UserInfo["email"]
	if ok {
		t[emailKey] = user.Email
	}
}
