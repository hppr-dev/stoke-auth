package web

import (
	"fmt"
	"time"
	"net/http"

	"stoke/internal/cfg"
	"stoke/internal/adm"
)

type Server struct {
	Config cfg.Server
	Server *http.Server
	handlers http.Handler
}

func (s *Server) Init() error {
	s.Server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServerFS(adm.Pages)))
	http.HandleFunc("/api", RootApiHandler)

	return nil
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}
