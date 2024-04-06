package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"stoke/client/stoke"
	"stoke/internal/admin"
	"stoke/internal/cfg"
	"stoke/internal/key"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(ctx context.Context) *http.Server {
	config := cfg.Ctx(ctx).Server
	issuer := key.IssuerFromCtx(ctx)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Address, config.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        InjectContext(ctx, LogHTTP(TraceHTTP(mux))),
	}

	// Static files
	mux.Handle("/admin/", http.StripPrefix("/admin/", http.FileServerFS(admin.Pages)))

	// TODO restrict methods/origins
	mux.Handle(
		"/api/login",
		AllowAllMethods(
			LoginApiHandler{},
		),
	)

	mux.Handle(
		"/api/pkeys",
		AllowAllMethods(
			PkeyApiHandler{},
		),
	)

	mux.Handle(
		"/api/refresh",
		AllowAllMethods(
			stoke.Auth(
				RefreshApiHandler{},
				issuer,
				stoke.Token().ForAccess(),
			),
		),
	)

	mux.Handle(
		"/api/admin_users",
		AllowAllMethods(
			stoke.Auth(
				UserHandler{},
				issuer,
				stoke.Token().Requires("srol", "spr").ForAccess(),
			),
		),
	)

	mux.Handle(
		"/api/admin/",
		AllowAllMethods(
			stoke.Auth(
				NewEntityAPIHandler("/api/admin/", ctx),
				issuer,
				stoke.Token().Requires("srol", "spr").ForAccess(),
			),
		),
	)

	mux.Handle(
		"/metrics",
		AllowAllMethods(
				promhttp.Handler(),
		),
	)

	return server
}
