package web

import (
	"errors"
	"fmt"
	"net/http"
	"stoke/internal/ctx"

	"github.com/go-faster/jx"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshApiHandler struct {
	Context *ctx.Context
}

// Request takes refresh token only. Must be authenticated by a valid token
func (r RefreshApiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
	newToken, newRefresh, err := r.Context.Issuer.RefreshToken(token, refresh, r.Context.Config.Tokens.TokenDuration)
	if err != nil {
		logger.Debug().Err(err).Str("refresh", refresh).Msg("Failed to refresh token")
		BadRequest.Write(res)
		return
	}

	res.Write([]byte(fmt.Sprintf("{\"token\":\"%s\",\"refresh\":\"%s\"}", newToken, newRefresh)))
}
