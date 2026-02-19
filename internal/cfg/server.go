package cfg

import (
	"context"
	"net/http"
)

type Server struct {
	// Address to serve the application on
	Address       string `json:"address"`
	// Port to serve the application on
	Port          int    `json:"port"`
	// Request Timeout in milliseconds
	Timeout       int    `json:"timeout"`
	// Base Path for web assets and api, useful for hosting behind a proxy
	BasePath   string `json:"base_path"`

	// Private TLS Key
	TLSPrivateKey string `json:"tls_private_key"`
	// Public TLS Cert
	TLSPublicCert string `json:"tls_public_cert"`

	// Allowed Hosts
	AllowedHosts []string `json:"allowed_hosts"`

	// Disable the admin UI
	DisableAdmin bool     `json:"disable_admin"`
}

func (s Server) WithContext(ctx context.Context) context.Context {
	mux := http.NewServeMux()
	return context.WithValue(ctx, serveMuxCtxKey, mux)
}

func MuxFromContext(ctx context.Context) *http.ServeMux {
	return ctx.Value(serveMuxCtxKey).(*http.ServeMux)
}
