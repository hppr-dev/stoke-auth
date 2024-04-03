package web

import (
	"fmt"
	"net/http"
	"time"

	"stoke/client/stoke"
	"stoke/internal/admin"
	"stoke/internal/ctx"
)

type Server struct {
	Context *ctx.Context
	Server *http.Server
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
				LoginApiHandler{ Context: s.Context },
			),
		),
	)

	http.Handle(
		"/api/pkeys",
		AllowAllMethods(
			LogHTTP(
				PkeyApiHandler{ Context: s.Context },
			),
		),
	)

	http.Handle(
		"/api/refresh",
		AllowAllMethods(
			LogHTTP(
				stoke.Auth(
					RefreshApiHandler{ Context: s.Context },
					s.Context.Issuer,
					stoke.Token().ForAccess(),
				),
			),
		),
	)

	http.Handle(
		"/api/admin_users",
		AllowAllMethods(
			LogHTTP(
				stoke.Auth(
					UserHandler{ Context: s.Context },
					s.Context.Issuer,
					stoke.Token().Requires("srol", "spr").ForAccess(),
				),
			),
		),
	)

	http.Handle(
		"/api/admin/",
		AllowAllMethods(
			LogHTTP(
				stoke.Auth(
					NewEntityAPIHandler("/api/admin/", s.Context),
					s.Context.Issuer,
					stoke.Token().Requires("srol", "spr").ForAccess(),
				),
			),
		),
	)

	return nil
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
