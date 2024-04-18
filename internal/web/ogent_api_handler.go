package web

import (
	"context"
	"net/http"
	"stoke/client/stoke"
	"stoke/internal/ent"
	"stoke/internal/ent/ogent"
	"stoke/internal/key"
	"stoke/internal/usr"

	"github.com/rs/zerolog"
)

func NewEntityAPIHandler(prefix string, ctx context.Context) http.Handler {
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[0 : len(prefix)-1]
	}

	eHandler := &entityHandler{}
	eHandler.Init(ctx)

	sHandler := &secHandler{
		TokenHandler: stoke.NewTokenHandler(
			key.IssuerFromCtx(ctx),
			stoke.WithToken().Requires("srol", "spr").ForAccess(),
		),
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
	*stoke.TokenHandler
}

// HandleToken implements ogent.SecurityHandler.
func (s *secHandler) HandleToken(ctx context.Context, operationName string, t ogent.Token) (context.Context, error) {
	return s.TokenHandler.InjectToken(t.GetToken(), ctx)
}

type entityHandler struct {
	*ogent.OgentHandler
}

func (h *entityHandler) Init(ctx context.Context) {
	h.OgentHandler = ogent.NewOgentHandler(ent.FromContext(ctx))
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
	if err := usr.ProviderFromCtx(ctx).AddUser(req.Fname, req.Lname, req.Email, req.Username, req.Password, false, ctx); err != nil {
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
