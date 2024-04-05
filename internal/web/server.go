package web

import (
	"fmt"
	"net/http"
	"time"

	"stoke/client/stoke"
	"stoke/internal/admin"
	"stoke/internal/ctx"
	"stoke/internal/tel"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Context *ctx.Context
	Server *http.Server
	OTEL   *tel.OTEL
}

func (s *Server) Init() error {
	serverAddr := fmt.Sprintf("%s:%d", s.Context.Config.Server.Address, s.Context.Config.Server.Port)
	s.Server = &http.Server{
		Addr:           serverAddr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Static files
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServerFS(admin.Pages)))

	// TODO restrict methods/origins
	// Register handlers
	http.Handle(
		"/api/login",
		AllowAllMethods(
			LogHTTP(
				WithSpan("LOGIN",
					LoginApiHandler{ Context: s.Context },
				),
			),
		),
	)

	http.Handle(
		"/api/pkeys",
		AllowAllMethods(
			LogHTTP(
				WithSpan("PKEYS",
					PkeyApiHandler{ Context: s.Context },
				),
			),
		),
	)

	http.Handle(
		"/api/refresh",
		AllowAllMethods(
			LogHTTP(
				WithSpan("REFRESH",
					stoke.Auth(
						RefreshApiHandler{ Context: s.Context },
						s.Context.Issuer,
						stoke.Token().ForAccess(),
					),
				),
			),
		),
	)

	http.Handle(
		"/api/admin_users",
		AllowAllMethods(
			LogHTTP(
				WithSpan("ADMINUSERS",
					stoke.Auth(
						UserHandler{ Context: s.Context },
						s.Context.Issuer,
						stoke.Token().Requires("srol", "spr").ForAccess(),
					),
				),
			),
		),
	)

	http.Handle(
		"/api/admin/",
		AllowAllMethods(
			LogHTTP(
				WithSpan("ADMINAPI",
					stoke.Auth(
						NewEntityAPIHandler("/api/admin/", s.Context, s.OTEL),
						s.Context.Issuer,
						stoke.Token().Requires("srol", "spr").ForAccess(),
					),
				),
			),
		),
	)

	http.Handle(
		"/metrics",
		AllowAllMethods(
			LogHTTP(
				promhttp.Handler(),
			),
		),
	)

	return nil
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
