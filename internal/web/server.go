package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"stoke/internal/admin"
	"stoke/internal/cfg"
	"stoke/internal/key"

	"hppr.dev/stoke"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type debugLogger struct {
	logger zerolog.Logger
}

func (d debugLogger) Write(p []byte) (int, error) {
	d.logger.Debug().Msg(strings.TrimSuffix(string(p), "\n"))
	return len(p), nil
}

func NewServer(ctx context.Context) *http.Server {
	logger := zerolog.Ctx(ctx)
	config := cfg.Ctx(ctx).Server
	logFile := cfg.Ctx(ctx).Logging.LogFile
	telConf := cfg.Ctx(ctx).Telemetry
	issuer := key.IssuerFromCtx(ctx)

	mux := http.NewServeMux()

	dLogger := debugLogger{ logger: logger.With().Str("component", "http.Server").Logger() }

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Address, config.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        InjectContext(ctx, TraceHTTP(LogHTTP(mux))),
		ErrorLog:       log.New(dLogger, "", 0),
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

	if !config.DisableAdmin {
		// Static files
		mux.Handle("/admin/", http.StripPrefix("/admin/", http.FileServerFS(admin.Pages)))
	}

	allowedHosts := strings.Join(config.AllowedHosts,",")

	mux.Handle(
		"/api/",
		ConfigureCORS(
			"GET,POST,PATCH,DELETE,OPTIONS",
			allowedHosts,
			NewEntityAPIHandler("/api/", ctx),
		),
	)


	if !telConf.DisableMonitoring {
		logsHandler := func(res http.ResponseWriter, req *http.Request) {
			logger := zerolog.Ctx(ctx)
			f, err := os.Open(logFile)
			if err != nil {
				logger.Error().Err(err).Msg("Could not open log file")
			}

			// Only return the last 6144 bytes of the logfile
			newOffset, err := f.Seek(-6144, 2)
			if err != nil {
				logger.Debug().Err(err).Int64("newOffset", newOffset).Msg("Could not seek log file")
			}

			content, err := io.ReadAll(f)
			if err != nil {
				logger.Error().Err(err).Msg("Could not read log file")
			}

			res.WriteHeader(http.StatusOK)
			res.Write(content)
		}

		if telConf.RequirePrometheusAuthentication {
			mux.Handle(
				"/metrics",
				ConfigureCORS(
					"GET,OPTIONS",
					allowedHosts,
					stoke.Auth(
						promhttp.Handler(),
						issuer,
						stoke.RequireToken().WithClaim("stk", "m").Or(
							stoke.RequireToken().WithClaim("stk", "s"),
						).Or(
							stoke.RequireToken().WithClaim("stk", "S"),
						),
					),
				),
			)
			mux.Handle(
				"/metrics/logs",
				ConfigureCORS(
					"GET,OPTIONS",
					allowedHosts,
					stoke.AuthFunc(
						logsHandler,
						issuer,
						stoke.RequireToken().WithClaim("stk", "m").Or(
							stoke.RequireToken().WithClaim("stk", "s"),
						).Or(
							stoke.RequireToken().WithClaim("stk", "S"),
						),
					),
				),
			)

		} else {
			mux.Handle(
				"/metrics",
				ConfigureCORS(
					"GET,OPTIONS",
					allowedHosts,
					promhttp.Handler(),
				),
			)
			mux.Handle(
				"/metrics/logs",
				ConfigureCORSFunc(
					"GET,OPTIONS",
					allowedHosts,
					logsHandler,
				),
			)

		}

	}

	return server
}
