package web

import (
	"context"
	"net/http"
	"stoke/internal/ent"
	"stoke/internal/ent/ogent"
	"stoke/internal/usr"

	"github.com/rs/zerolog"
)

func NewEntityAPIHandler(prefix string, ctx context.Context) http.Handler {
	if len(prefix) > 0 && prefix[len(prefix) - 1] == '/' {
		prefix = prefix[0:len(prefix) - 1]
	}

	eHandler := &entityHandler{}
	eHandler.Init(ctx)

	hdlr, err := ogent.NewServer(
		eHandler,
		ogent.WithPathPrefix(prefix),
	)
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("An error occured creating an entity handler")
		return nil
	}

	// Unwrap ServeHTTP to be have tracing spans use context from our custom middleware
	return http.HandlerFunc(hdlr.ServeHTTP)
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

func (h *entityHandler) CreateLocalUser(ctx context.Context, req ogent.OptCreateLocalUserReq) (ogent.CreateLocalUserRes, error) {
	v := req.Value

	if err := usr.ProviderFromCtx(ctx).AddUser(v.Fname, v.Lname, v.Email, v.Username, v.Password, false, ctx) ; err != nil {
		return &ogent.CreateLocalUserBadRequest{
			Message: err.Error(),
		}, nil
	}

	return &ogent.CreateLocalUserOK{}, nil
}

func (h *entityHandler) UpdateLocalUserPassword(ctx context.Context, req ogent.OptUpdateLocalUserPasswordReq) (ogent.UpdateLocalUserPasswordRes, error) {
	v := req.Value

	if err := usr.ProviderFromCtx(ctx).UpdateUserPassword(v.Username, v.OldPassword.Value, v.NewPassword, v.Force.Or(false), ctx); err != nil {
		return &ogent.UpdateLocalUserPasswordBadRequest{
			Message: err.Error(),
		}, nil
	}

	return &ogent.UpdateLocalUserPasswordOK{}, nil
}
