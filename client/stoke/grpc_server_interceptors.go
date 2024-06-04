package stoke

import (
	"context"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor to authenticate unary grpc requests.
// Use NewUnaryTokenInterceptor to create.
type UnaryTokenInterceptor struct {
	preInterceptors []*UnaryTokenInterceptor
	methodFilter *regexp.Regexp
	shortcut bool
	skip bool
	*TokenHandler
}

// Create a new unary token interceptor
// 
// Example Usage:
//   s := grpc.NewServer(NewUnaryTokenInterceptor(
//     stoke.NewPerRequestPublicKeyStore("http://stoke:8080/api/pkeys", context.Background()),
//     stoke.RequireToken(),
//   ))
func NewUnaryTokenInterceptor(store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) *UnaryTokenInterceptor {
	return &UnaryTokenInterceptor{
		TokenHandler: NewTokenHandler(store, claims, parserOpts...),
	}
}

// This creates a UnaryToken interceptor that authenticates claims based on a map of grpc method names (regexps) to claims
// Note: processing order is not guarenteed, i.e. map keys should match only one route
//
// Example:
//   routeInt := NewUnaryInterceptorRouteMap(store, map[string]*Claims{
//     "/service/one" : stoke.RequireToken().WithClaim("role", "rep"),
//     "/service/two" : stoke.RequireToken().WithClaim("role", "user"),
//   })
func NewUnaryInterceptorRouteMap(store PublicKeyStore, claimMap map[string]*Claims, parserOpts ...jwt.ParserOption) *UnaryTokenInterceptor {
	root := NewUnaryTokenInterceptor(store, RequireToken(), parserOpts...).Shortcut().SkipHandler()

	for method, claims := range claimMap {
		root.Intercept(NewUnaryTokenInterceptor(store, claims, parserOpts...).Method(method))
	}

	return root
}

// Shortcut processing when an interceptor returns a non-nil response
func (uti *UnaryTokenInterceptor) Shortcut() *UnaryTokenInterceptor {
	uti.shortcut = true
	return uti
}

// Skip calling the grpc unary handler and return a nil response and error when auth succeedes.
// Note that the resulting context in the unary handler will not include a token.
func (uti *UnaryTokenInterceptor) SkipHandler() *UnaryTokenInterceptor {
	uti.skip = true
	return uti
}

// Set the method name by regex.
//
// If the method name passed in the interceptor info does not match the given regex, processing stops and returns a nil response and error.
// Does not filter on method name if not called
func (uti *UnaryTokenInterceptor) Method(methodRegex string) *UnaryTokenInterceptor {
	uti.methodFilter = regexp.MustCompile(methodRegex)
	return uti
}

// Places an interceptor to be processed before this interceptor.
// Allows the use of multiple unary token interceptors.
//
// If Shortcut has been set, then any non-nil response will result in skipping this interceptor's processing.
func (uti *UnaryTokenInterceptor) Intercept(usi *UnaryTokenInterceptor) *UnaryTokenInterceptor {
	uti.preInterceptors = append(uti.preInterceptors, usi)
	return uti
}

// Converts this UnaryTokenInterceptor into a grpc server option
func (uti *UnaryTokenInterceptor) Opt() grpc.ServerOption {
	return grpc.UnaryInterceptor(uti.innerIntercept)
}

// Intercepts all requests, authenticates the token and injects the context
func (uti *UnaryTokenInterceptor) innerIntercept(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	for _, inter := range uti.preInterceptors {
		res, err := inter.innerIntercept(ctx, req, info, handler)
		if err != nil || ( uti.shortcut && res != nil ){
			return res, err
		}
	}

	if uti.methodFilter != nil && !uti.methodFilter.MatchString(info.FullMethod) {
		return nil, nil
	}

	token, err := verifyMetadata(ctx)
	if err != nil {
		return nil, err
	}

	fwCtx, err := uti.InjectToken(token, ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid Token")
	}

	if uti.skip {
		return nil, nil
	}

	resp, err := handler(fwCtx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	return resp, nil
}

// Interceptor to authenticate streaming grpc requests.
// Use NewStreamTokenInterceptor to create.
type StreamTokenInterceptor struct {
	preInterceptors []*StreamTokenInterceptor
	methodFilter *regexp.Regexp
	shortcut bool
	skip bool
	*TokenHandler
}

// Create a new unary token interceptor
// 
// Example Usage:
//   s := grpc.NewServer(NewStreamTokenInterceptor(
//     stoke.NewPerRequestPublicKeyStore("http://stoke:8080/api/pkeys", context.Background()),
//     stoke.RequireToken(),
//   ))
func NewStreamTokenInterceptor(store PublicKeyStore, claims *Claims, parserOpts ...jwt.ParserOption) *StreamTokenInterceptor {
	return &StreamTokenInterceptor{
		TokenHandler: NewTokenHandler(store, claims, parserOpts...),
	}
}

// Shortcut processing when an interceptor returns a non-nil response
func (sti *StreamTokenInterceptor) Shortcut() *StreamTokenInterceptor {
	sti.shortcut = true
	return sti
}

// Skip calling the grpc unary handler and return a nil response and error when auth succeedes.
// Note that the resulting context in the unary handler will not include a token.
func (sti *StreamTokenInterceptor) SkipHandler() *StreamTokenInterceptor {
	sti.skip = true
	return sti
}

// Set the method name by regex.
//
// If the method name passed in the interceptor info does not match the given regex, processing stops and returns a nil response and error.
// Does not filter on method name if not called
func (sti *StreamTokenInterceptor) Method(methodRegex string) *StreamTokenInterceptor {
	sti.methodFilter = regexp.MustCompile(methodRegex)
	return sti
}

// Places a stream token interceptor to be processed before this interceptor.
// Allows the use of multiple unary interceptors.
//
// If Shortcut has been set, then any non-nil response will result in skipping this interceptor's processing.
func (sti *StreamTokenInterceptor) Intercept(ssi *StreamTokenInterceptor) *StreamTokenInterceptor {
	sti.preInterceptors = append(sti.preInterceptors, ssi)
	return sti
}

// Converts this StreamTokenInterceptor into a grpc server option
func (sti *StreamTokenInterceptor) Opt() grpc.ServerOption {
	return grpc.StreamInterceptor(sti.innerIntercept)
}

// Intercepts the initial request, authenticates the token and injects the context
func (sti *StreamTokenInterceptor) innerIntercept(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	for _, inter := range sti.preInterceptors {
		err := inter.innerIntercept(srv, ss, info, handler)
		if err != nil || sti.shortcut {
			return err
		}
	}

	if sti.methodFilter != nil && !sti.methodFilter.MatchString(info.FullMethod) {
		return nil
	}

	ctx := ss.Context()

	token, err := verifyMetadata(ctx)
	if err != nil {
		return err
	}

	fwCtx, err := sti.InjectToken(token, ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Invalid Token")
	}

	if sti.skip {
		return nil
	}

	return handler(srv, &streamWrap{
		ctx: fwCtx,
		ServerStream: ss,
	})
}

type streamWrap struct {
	ctx context.Context
	grpc.ServerStream
}

func (sw streamWrap) Context() context.Context {
	return sw.ctx
}


func verifyMetadata(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "No request metadata")
	}

	auth, ok := meta["authorization"]
	if !ok || len(auth) != 1 || !strings.HasPrefix(auth[0], "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "Missing or malformed authorization metadata: %v", meta)
	}

	return strings.TrimPrefix(auth[0], "Bearer "), nil
}
