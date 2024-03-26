package web

import (
	"fmt"
	"net/http"
	"time"

	"stoke/internal/adm"
	"stoke/internal/ctx"
)

type Server struct {
	Context *ctx.Context
	Server *http.Server
}

func (s *Server) Init() error {
	s.Server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.Context.Config.Server.Address, s.Context.Config.Server.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServerFS(adm.Pages)))
	http.Handle("/api/login", LoginApiHandler{ Context: s.Context } )
	http.Handle("/api/pkeys", PkeyApiHandler{ Context: s.Context } )

	return nil
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
