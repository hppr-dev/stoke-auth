package web

import "net/http"

func ConfigureCORS( methods, origins string, h http.Handler ) corsWrapper {
	return corsWrapper{
		inner : h.ServeHTTP,
		allowedMethods : methods,
		allowedOrigins : origins,
	}
}

func ConfigureCORSFunc( methods, origins string, h http.HandlerFunc ) corsWrapper {
	return corsWrapper{
		inner : h,
		allowedMethods : methods,
		allowedOrigins : origins,
	}
}

func AllowAllMethods(h http.Handler) corsWrapper {
	return ConfigureCORS( "GET, POST, PUT, DELETE, PATCH, OPTIONS", "*", h)
}

func AllowAllMethodsFunc(h http.HandlerFunc) corsWrapper {
	return ConfigureCORSFunc( "GET, POST, PUT, DELETE, PATCH, OPTIONS", "*", h)
}

type corsWrapper struct {
	inner http.HandlerFunc
	allowedMethods string
	allowedOrigins string
}

func (w corsWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Access-Control-Allow-Origin", w.allowedOrigins)
	if req.Method == http.MethodOptions {
		headers := res.Header()
		headers.Add("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
		headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization")
		headers.Add("Access-Control-Allow-Methods", w.allowedMethods)
		res.WriteHeader(http.StatusNoContent)
		return
	}
	w.inner(res, req)
}
