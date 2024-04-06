package web

import (
	"errors"
	"fmt"
	"net/http"
	"stoke/internal/cfg"
	"stoke/internal/key"

	"github.com/go-faster/jx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type RefreshApiHandler struct {}

// Request takes refresh token only. Must be authenticated by a valid token
func (r RefreshApiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := zerolog.Ctx(ctx)

	if req.Method != http.MethodPost {
		MethodNotAllowed.Write(res)
		return
	}

	var refresh string
	decoder := jx.Decode(req.Body, 256)
	err := decoder.Obj(func (d *jx.Decoder, key string) error {
		var err error
		switch key {
		case "refresh":
			refresh, err = d.Str()
		default:
			return errors.New("Bad Request")
		}
		return err
	})

	if err != nil || refresh == "" {
		logger.Debug().Err(err).Str("refresh", refresh).Msg("Missing body parameters")
		BadRequest.Write(res)
		return
	}

	token := req.Context().Value("jwt.Token").(*jwt.Token)
	newToken, newRefresh, err := key.IssuerFromCtx(ctx).RefreshToken(token, refresh, cfg.Ctx(ctx).Tokens.TokenDuration)
	if err != nil {
		logger.Debug().Err(err).Str("refresh", refresh).Msg("Failed to refresh token")
		BadRequest.Write(res)
		return
	}

	res.Write([]byte(fmt.Sprintf("{\"token\":\"%s\",\"refresh\":\"%s\"}", newToken, newRefresh)))
}
