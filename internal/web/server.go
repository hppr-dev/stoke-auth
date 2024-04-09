package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"stoke/client/stoke"
	"stoke/internal/admin"
	"stoke/internal/cfg"
	"stoke/internal/key"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func NewServer(ctx context.Context) *http.Server {
	logger := zerolog.Ctx(ctx)
	config := cfg.Ctx(ctx).Server
	issuer := key.IssuerFromCtx(ctx)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Address, config.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        InjectContext(ctx, TraceHTTP(LogHTTP(mux))),
	}

	if config.TLSPrivateKey != "" && config.TLSPublicCert != "" {
		cert, err := tls.LoadX509KeyPair(config.TLSPublicCert, config.TLSPrivateKey)
		if err != nil {
			logger.Fatal().
				Err(err).
				Str("privateKey", config.TLSPrivateKey).
				Str("publicKey", config.TLSPublicCert).
				Msg("Could not load TLS certificates")
		}
		server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		server.TLSConfig = &tls.Config{
			Certificates:             []tls.Certificate{ cert },
			MinVersion:               tls.VersionTLS12,
		}
	} else {
		logger.Error().
			Str("privateKey", config.TLSPrivateKey).
			Str("publicKey", config.TLSPublicCert).
			Str("error", "tls_private_key and/or tls_public_key not set.").
			Str("advice", "Enable TLS in production environments").
			Msg("Failed to initialize TLS.")
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
				stoke.WithToken().ForAccess(),
			),
		),
	)

	mux.Handle(
		"/api/admin_users",
		AllowAllMethods(
			stoke.Auth(
				UserHandler{},
				issuer,
				stoke.WithToken().Requires("srol", "spr").ForAccess(),
			),
		),
	)

	mux.Handle(
		"/api/admin/",
		AllowAllMethods(
			stoke.Auth(
				NewEntityAPIHandler("/api/admin/", ctx),
				issuer,
				stoke.WithToken().Requires("srol", "spr").ForAccess(),
			),
		),
	)

	mux.Handle(
		"/metrics",
		AllowAllMethods(
			stoke.Auth(
				promhttp.Handler(),
				issuer,
				stoke.WithToken().Requires("srol", "spr").ForAccess(),
			),
		),
	)

	return server
}
