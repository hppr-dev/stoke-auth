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

type baseTokenInterceptor struct {
	methodFilter *regexp.Regexp
	shortcut bool
	skip bool
	*TokenHandler
}

// Shortcut processing when an interceptor returns a non-nil response
func (bti *baseTokenInterceptor) Shortcut() *baseTokenInterceptor {
	bti.shortcut = true
	return bti
}

// Skip calling the grpc unary handler and return a nil response and error when auth succeedes.
// Note that the resulting context in the unary handler will not include a token.
func (bti *baseTokenInterceptor) SkipHandler() *baseTokenInterceptor {
	bti.skip = true
	return bti
}

// Set the method name by regex.
//
// If the method name passed in the interceptor info does not match the given regex, processing stops and returns a nil response and error.
// Does not filter on method name if not called
func (bti *baseTokenInterceptor) Method(methodRegex string) *baseTokenInterceptor {
	bti.methodFilter = regexp.MustCompile(methodRegex)
	return bti
}

// Interceptor to authenticate unary grpc requests.
// Use NewUnaryTokenInterceptor to create.
type UnaryTokenInterceptor struct {
	preInterceptors []grpc.UnaryServerInterceptor
	baseTokenInterceptor
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
		baseTokenInterceptor:  baseTokenInterceptor {
			TokenHandler: NewTokenHandler(store, claims, parserOpts...),
		},
	}
}

// Places an interceptor to be processed before this interceptor.
// Allows the use of multiple unary interceptors.
//
// If Shortcut has been set, then any non-nil response will result in skipping this interceptor's processing.
func (uti *UnaryTokenInterceptor) Intercept(usi grpc.UnaryServerInterceptor) *UnaryTokenInterceptor {
	uti.preInterceptors = append(uti.preInterceptors, usi)
	return uti
}

// Converts this UnaryTokenInterceptor into a grpc server option
func (uti *UnaryTokenInterceptor) Opt() grpc.ServerOption {
	return grpc.UnaryInterceptor(uti.innerIntercept)
}

// Intercepts all requests, authenticates the token and injects the context
func (uti *UnaryTokenInterceptor) innerIntercept(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if uti.methodFilter != nil && uti.methodFilter.MatchString(info.FullMethod) {
		return nil, nil
	}

	for _, inter := range uti.preInterceptors {
		res, err := inter(ctx, req, info, handler)
		if err != nil || ( uti.shortcut && res != nil ){
			return res, err
		}
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
	preInterceptors []grpc.StreamServerInterceptor
	baseTokenInterceptor
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
		baseTokenInterceptor:  baseTokenInterceptor {
			TokenHandler: NewTokenHandler(store, claims, parserOpts...),
		},
	}
}

// Places an interceptor to be processed before this interceptor.
// Allows the use of multiple unary interceptors.
//
// If Shortcut has been set, then any non-nil response will result in skipping this interceptor's processing.
func (sti *StreamTokenInterceptor) Intercept(ssi grpc.StreamServerInterceptor) *StreamTokenInterceptor {
	sti.preInterceptors = append(sti.preInterceptors, ssi)
	return sti
}

// Converts this StreamTokenInterceptor into a grpc server option
func (sti *StreamTokenInterceptor) Opt() grpc.ServerOption {
	return grpc.StreamInterceptor(sti.innerIntercept)
}

// Intercepts the initial request, authenticates the token and injects the context
func (sti *StreamTokenInterceptor) innerIntercept(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if sti.methodFilter != nil && sti.methodFilter.MatchString(info.FullMethod) {
		return nil
	}

	for _, inter := range sti.preInterceptors {
		err := inter(srv, ss, info, handler)
		if err != nil || sti.shortcut {
			return err
		}
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

	auth, ok := meta["Authorization"]
	if !ok || len(auth) != 1 || !strings.HasPrefix(auth[0], "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "Missing or malformed Authorization metadata")
	}

	return strings.TrimPrefix(auth[0], "Bearer "), nil
}
