package web

import (
	"net/http"
	"stoke/internal/key"

	"github.com/rs/zerolog"
)

type PkeyApiHandler struct {}

func (p PkeyApiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	b, err := key.IssuerFromCtx(ctx).PublicKeys()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Could not get public keys")
		InternalServerError.Write(res)
		return
	}
	res.Write(b)
}
