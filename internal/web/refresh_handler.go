package web

import (
	"context"
	"stoke/internal/cfg"
	"stoke/internal/ent/ogent"
	"stoke/internal/key"
	"stoke/internal/tel"

	"hppr.dev/stoke"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

// Refresh implements ogent.Handler.
func (h *entityHandler) Refresh(ctx context.Context, req *ogent.RefreshReq) (ogent.RefreshRes, error) {
	logger := zerolog.Ctx(ctx)

	ctx, span := tel.GetTracer().Start(ctx, "RefreshHandler")
	defer span.End()

	newToken, newRefresh, err := key.IssuerFromCtx(ctx).RefreshToken(stoke.Token(ctx), req.Refresh, cfg.Ctx(ctx).Tokens.TokenDuration, ctx)
	if err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("refresh", req.Refresh).
			Msg("Failed to refresh token")
		return &ogent.RefreshUnauthorized{}, nil
	}

	return &ogent.RefreshOK{
		Token:   newToken,
		Refresh: newRefresh,
	}, nil
}
