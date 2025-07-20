package middleware

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
)

func AuthUnaryInterceptor(maker token.Maker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata missing")
		}
		auths := md.Get("authorization")
		if len(auths) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization header missing")
		}
		fields := strings.Fields(auths[0])
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			return nil, status.Errorf(codes.Unauthenticated, "bad auth header")
		}
		payload, err := maker.VerifyToken(fields[1])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}
		newCtx := context.WithValue(ctx, "payload", payload)
		return handler(newCtx, req)
	}
}
