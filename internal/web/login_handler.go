package web

import (
	"context"
	"fmt"
	"slices"
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/ent/ogent"
	"stoke/internal/key"
	"stoke/internal/tel"
	"stoke/internal/usr"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"hppr.dev/stoke"
)

type LoginApiHandler struct {}

//   1. Retrieves claims from user provider using username and password
//   2. Verifies required user claims match the requested values. If a the required value is "", the key is required but no further verification is needed.
//   3. Removes any claims that do not match the claim_filter, if given
//   3. Issues a token
// Schema definition in internal/schema/openapi/login.go and internal/ent/openapi.json (operation id login)
func (h *entityHandler) Login(ctx context.Context, req *ogent.LoginReq) (ogent.LoginRes, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "Login").
		Str("username", req.Username).
		Strs("filter_claims", req.FilterClaims).
		Logger()


	ctx, span := tel.GetTracer().Start(ctx, "LoginHandler")
	defer span.End()

	user, pvClaims, err := usr.ProviderFromCtx(ctx).GetUserClaims(req.Username, req.Password, ctx)
	if err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to get claims from provider")
		return &ogent.LoginUnauthorized{}, nil
	}
	
	tokenMap := make(map[string]string)

	matchedOne := len(req.RequiredClaims) == 0
	for _, pvClaim := range pvClaims {
		if !matchedOne {
			for _, claimReq := range req.RequiredClaims {
				if value, exist := claimReq[pvClaim.ShortName]; exist {
					matchedOne = value == "" || pvClaim.Value == value
				}
			}
		}

		if len(req.FilterClaims) == 0 || slices.Contains(req.FilterClaims, pvClaim.ShortName) {
			if value, exists := tokenMap[pvClaim.ShortName]; exists {
				tokenMap[pvClaim.ShortName] = value + "," + pvClaim.Value
			} else {
				tokenMap[pvClaim.ShortName] = pvClaim.Value
			}
		}
	}

	if !matchedOne {
		logger.Debug().
			Interface("provider_claims", pvClaims).
			Interface("required_claimes", req.RequiredClaims).
			Msg("User did not have a required claimset.")
		return &ogent.LoginUnauthorized{}, nil
	}

	populateUserInfo(cfg.Ctx(ctx), user, tokenMap)

	token, refresh, err := key.IssuerFromCtx(ctx).IssueToken(&stoke.Claims{
		StokeClaims : tokenMap,
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
