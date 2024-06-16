package stoke

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Requires authentication for access to a handler
// Wraps the given http.Handler to:
//	* Verify token signature
//	* Verify required claims
//	* Inject token into the request context.
//
// Example Usage:
//
//	type myHandler struct {}
//	func (myHandler) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {rs.Write([]byte(fmt.Sprintf("I got token: %s", stoke.Token(rq.Context()).Raw)))}
//
//	func ConfigureHttp() {
//		publicKeyStore := stoke.NewPerRequestPublicKeyStore()
//		http.Handle(
//			"/my-token",
//			stoke.Auth(myHandler{}, publicKeyStore, stoke.RequireToken()),
//		)
//	}
func Auth(handler http.Handler, store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler.ServeHTTP,
		TokenHandler: NewTokenHandler(store, claims, parserOpts...),
	}
}

// Requires authentication for access to a handler function
// Wraps the given http.HandlerFunc to:
//	* Verify token signature
//	* Verify required claims
//	* Inject token into the request context.
//
// Example Usage:
//
//	func ConfigureHttp() {
//		publicKeyStore := stoke.NewPerRequestPublicKeyStore()
//		http.Handle(
//			"/my-token",
//			stoke.AuthFunc(
//				func(rs http.ResponseWriter, rq *http.Request) { rs.Write([]byte(fmt.Sprintf("I got token: %s", stoke.Token(rq.Context()).Raw))) },
//				publicKeyStore,
//				stoke.RequireToken(),
//			),
//		)
//	}
func AuthFunc(handler http.HandlerFunc, store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) http.Handler {
	return authWrapper{
		inner:      handler,
		TokenHandler: NewTokenHandler(store, claims, parserOpts...),
	}
}

// Requires authentication for subpaths of a handler.
// Wraps the given http.Handler to:
//	* Verify token signature (rejects request if it doesn't match)
//	* Verify required claims by urls in a map[string]*stoke.Claims
//	* Inject token into the request context
//
// Uses a ServeMux internally to associate route handlers with claims, i.e. map keys are in [net/http.ServeMux] route form.
// The next handler will only be served if the route matches and the required claims are the expected value
// Set a route to nil to skip token auth for the given route
//
// Example Usage:
//
//	func ConfigureHttp() {
//		publicKeyStore := stoke.NewPerRequestPublicKeyStore()
//		http.Handle("/api/birds/", &BirdsHandler{}) // handle requests without worrying about auth
// 		...
//		rootHandler := stoke.AuthMap(http.DefaultServeMux, publicKeyStore,
//      map[string]*stoke.Claims{
//	      "/api/birds/": stoke.RequireToken().WithClaim("bird", "r"), // /api/birds/ requires bird:r
//	      "/api/bugs/":  stoke.RequireToken(),                        // /api/bugs/  requires any valid token
//	      "/api/bees/":  nil,                                         // /api/bees/  does not require a token at all
//			},
//    )
//	}
//
func AuthMap(next http.Handler, store PublicKeyStore, route2Claims map[string]*Claims, parserOpts ...jwt.ParserOption) http.Handler {
	mux := http.NewServeMux()
	for route, claims := range route2Claims {
		if claims == nil {
			mux.Handle(route, next)
		} else{
			mux.Handle(route, authWrapper{
				inner:      next.ServeHTTP,
				TokenHandler: NewTokenHandler(store, claims, parserOpts...),
			})
		}
	}
	return http.HandlerFunc(mux.ServeHTTP)
}

// Authentication wrapper that is used to require authentication
type authWrapper struct {
	inner http.HandlerFunc
	*TokenHandler
}

// Authentication wrapper ServeHTTP function.
// Returns 401 Unauthorized if:
//	* The token's signature cannot be verified
//	* Required claims do not match specified criteria
//
// Injects the token into the request context.
func (w authWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx, span := getTracer().Start(ctx, "AuthHandler.ServeHTTP")
	defer span.End()

	token := req.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte(`{"message" : "Token is required"}`))
		return
	}

	trimToken := strings.TrimPrefix(token, "Bearer ")

	ctx, err := w.InjectToken(trimToken, ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{"message" : "Unauthorized"}`))
		return
	}

	w.inner(res, req.WithContext(ctx))
}
