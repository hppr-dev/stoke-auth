package web

import (
	"context"
	"net/http"
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/ent/ogent"
	"stoke/internal/key"
	"stoke/internal/usr"

	"hppr.dev/stoke"
	"github.com/rs/zerolog"
)

func NewEntityAPIHandler(prefix string, ctx context.Context) http.Handler {
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[0 : len(prefix)-1]
	}

	eHandler := newEntityHandler(ctx)


	sHandler := &secHandler{
		PublicKeyStore : key.IssuerFromCtx(ctx),
	}

	hdlr, err := ogent.NewServer(
		eHandler,
		sHandler,
		ogent.WithPathPrefix(prefix),
	)
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("An error occured creating an entity handler")
		return nil
	}

	// Unwrap ServeHTTP to be have tracing spans use context from our custom middleware
	return http.HandlerFunc(hdlr.ServeHTTP)
}

type secHandler struct {
	stoke.PublicKeyStore
}

// HandleToken implements ogent.SecurityHandler.
func (s *secHandler) HandleToken(ctx context.Context, operationName string, t ogent.Token) (context.Context, error) {
	claims := stoke.RequireToken().WithClaim("stk", "S")
	switch operationName {
	case "ReadClaimGroup", "ReadGroupLink", "ReadGroupLinkClaimGroup", "ListClaimGroup", "ListClaimGroupClaims", "ListClaimGroupGroupLinks", "ListClaimGroupUsers", "ListGroupLink":
		claims.Or(stoke.RequireToken().WithClaimMatch("stk", "^[sgGU]$"))

	case "CreateClaimGroup", "DeleteClaimGroup", "UpdateClaimGroup", "CreateGroupLink", "DeleteGroupLink", "UpdateGroupLink":
		claims.Or(stoke.RequireToken().WithClaim("stk", "G"))

	case "ReadClaim", "ListClaim", "ListClaimClaimGroups" :
		claims.Or(stoke.RequireToken().WithClaimMatch("stk", "^[scCGU]$"))

	case "CreateClaim", "DeleteClaim", "UpdateClaim":
		claims.Or(stoke.RequireToken().WithClaim("stk", "C"))

	case "ReadUser", "ListUser", "ListUserClaimGroups":
		claims.Or(stoke.RequireToken().WithClaimMatch("stk", "^[suU]$"))

	case "DeleteUser", "UpdateUser", "CreateLocalUser", "UpdateLocalUserPassword":
		claims.Or(stoke.RequireToken().WithClaim("stk", "U"))

	case "Capabilities", "Totals":
		claims.Or(stoke.RequireToken().WithClaimMatch("stk", "^[sSuUgGcC]$"))

	case "Refresh":
		claims = stoke.RequireToken()

	//case "listPrivateKey", "readPrivateKey":
	// TODO either remove or incorperate
	}

	zerolog.Ctx(ctx).Debug().
		Str("operationName", operationName).
		Str("token", t.Token).
		Interface("reqClaims", claims).
		Msg("Checking token for operation")
	return stoke.NewTokenHandler(s.PublicKeyStore, claims).InjectToken(t.GetToken(), ctx)
}

type entityHandler struct {
	*ogent.OgentHandler
}

func newEntityHandler(ctx context.Context) *entityHandler {
	return &entityHandler{
		OgentHandler: ogent.NewOgentHandler(ent.FromContext(ctx)),
	}
}

// Capabilities implements ogent.Handler.
func (h *entityHandler) Capabilities(ctx context.Context) (*ogent.CapabilitiesOK, error) {
	var caps []string
	config := cfg.Ctx(ctx)

	if !config.Server.DisableAdmin {
		caps = append(caps, "admin")
	}
	
	if !config.Telemetry.DisableMonitoring {
		caps = append(caps, "monitoring")
	}

	return &ogent.CapabilitiesOK{
		Capabilities:   caps,
		BaseAdminPath: config.Server.BaseAdminPath,
	}, nil
}

func (h *entityHandler) Totals(ctx context.Context) (*ogent.TotalsOK, error) {
	client := ent.FromContext(ctx)
	userCount, err := client.User.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	claimCount, err := client.Claim.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	groupCount, err := client.ClaimGroup.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	return &ogent.TotalsOK{
		Users:       userCount,
		Claims:      claimCount,
		ClaimGroups: groupCount,
	}, nil
}

func (h *entityHandler) CreateLocalUser(ctx context.Context, req *ogent.CreateLocalUserReq) (ogent.CreateLocalUserRes, error) {
	if err := usr.ProviderFromCtx(ctx).AddUser(req.Fname, req.Lname, req.Email, req.Username, req.Password, ctx); err != nil {
		return &ogent.CreateLocalUserBadRequest{
			Message: err.Error(),
		}, nil
	}

	return &ogent.CreateLocalUserOK{}, nil
}

func (h *entityHandler) UpdateLocalUserPassword(ctx context.Context, req *ogent.UpdateLocalUserPasswordReq) (ogent.UpdateLocalUserPasswordRes, error) {
	if err := usr.ProviderFromCtx(ctx).UpdateUserPassword(req.Username, req.OldPassword.Value, req.NewPassword, req.Force.Or(false), ctx); err != nil {
		return &ogent.UpdateLocalUserPasswordBadRequest{
			Message: err.Error(),
		}, nil
	}

	return &ogent.UpdateLocalUserPasswordOK{}, nil
}
