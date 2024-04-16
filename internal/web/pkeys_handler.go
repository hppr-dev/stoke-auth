package web

import (
	"context"
	"stoke/internal/ent/ogent"
	"stoke/internal/key"

	"github.com/go-faster/jx"
	"github.com/rs/zerolog"
)

// Pkeys implements ogent.Handler.
func (h *entityHandler) Pkeys(ctx context.Context) (*ogent.PkeysOK, error) {
	b, err := key.IssuerFromCtx(ctx).PublicKeys(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Could not get public keys")
		return nil, err
	}
	res := &ogent.PkeysOK{}
	err = res.Decode(jx.DecodeBytes(b))
	return res, err
}
