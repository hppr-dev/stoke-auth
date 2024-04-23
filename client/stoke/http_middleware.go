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
//		publicKeyStore := stoke.DefaultPublicKeyStore()
//		publicKeyStore.Init(context.Background)
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
//		publicKeyStore := stoke.DefaultPublicKeyStore()
//		publicKeyStore.Init(context.Background)
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
