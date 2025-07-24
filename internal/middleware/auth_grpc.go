package middleware

// AuthUnaryInterceptor adds Paseto authentication to gRPC unary calls.
// All comments are in English.

import (
	"context"
	"strings"

	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthUnaryInterceptor verifies the Paseto token in the `authorization` metadata
// and injects the payload into the context. Expected header format:
//
//	authorization: Bearer <token>
func AuthUnaryInterceptor(maker token.Maker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata missing")
		}

		auths := md.Get("authorization")
		if len(auths) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header missing")
		}

		fields := strings.Fields(auths[0])
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid auth header")
		}

		payload, err := maker.VerifyToken(fields[1])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%v", err)
		}

		newCtx := context.WithValue(ctx, token.CtxKey, payload)
		return handler(newCtx, req)
	}
}
