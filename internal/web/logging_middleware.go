package web

import (
	"time"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func LogHTTPFunc(h http.HandlerFunc) http.Handler {
	return logWrapper{
		inner: h,
	}
}

func LogHTTP(h http.Handler) http.Handler {
	return logWrapper {
		inner: h.ServeHTTP,
	}
}

type logWrapper struct {
	inner http.HandlerFunc
}

func (w logWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		zerolog.Ctx(req.Context()).Debug().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("millis", duration).
				Send()
		},
	)(w.inner).ServeHTTP(res, req)
}

