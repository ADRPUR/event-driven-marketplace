package handler

// gRPC transport layer for the Auth domain.
// All comments are in English, as required.

import (
	"context"

	authv1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
)

// grpcServer implements authv1.AuthServiceServer.
type grpcServer struct {
	authv1.UnimplementedAuthServiceServer
	svc *service.AuthService
}

// NewGRPCServer returns a gRPC AuthServiceServer.
func NewGRPCServer(svc *service.AuthService) authv1.AuthServiceServer {
	return &grpcServer{svc: svc}
}

// ---------------- RPCs ----------------

func (s *grpcServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	at, rt, err := s.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &authv1.LoginResponse{AccessToken: at, RefreshToken: rt}, nil
}

func (s *grpcServer) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	at, err := s.svc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &authv1.RefreshResponse{AccessToken: at}, nil
}

func (s *grpcServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := s.svc.Logout(ctx, req.RefreshToken); err != nil {
		return nil, err
	}
	return &authv1.LogoutResponse{Success: true}, nil
}

// ---------------- Helpers ----------------

// PayloadFromCtx is a convenience accessor.
func PayloadFromCtx(ctx context.Context) *token.Payload {
	p, _ := ctx.Value(token.CtxKey).(*token.Payload)
	return p
}
