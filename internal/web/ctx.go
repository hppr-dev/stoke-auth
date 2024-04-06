package web

import (
	"context"
	"net/http"
	"stoke/internal/cfg"
	"time"
)


func InjectContext(rootCtx context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(
		func (res http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(rootCtx, time.Millisecond * time.Duration(cfg.Ctx(rootCtx).Server.Timeout))
			defer cancel()
			h.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
