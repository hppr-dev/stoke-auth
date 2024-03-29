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

	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServerFS(admin.Pages)))

	http.Handle("/api/login", LogHTTP(LoginApiHandler{ Context: s.Context }) )
	http.Handle("/api/pkeys", LogHTTP(PkeyApiHandler{ Context: s.Context }) )

	entityPrefix := "/api/admin/"
	http.Handle(entityPrefix,
		LogHTTP(
			stoke.WithClaims(
				NewEntityAPIHandler(entityPrefix, s.Context),
				s.Context.Issuer,
				stoke.Claims().Require("srol", "spr").Validator(),
			),
		),
	)

	return nil
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
