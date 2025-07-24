package handler

// gRPC transport layer for the Auth domain.
// All comments are in English, as required.

import (
	"context"
	"time"

	auth1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// grpcServer implements AuthServiceServer
type grpcServer struct {
	auth1.UnimplementedAuthServiceServer
	svc *service.Service
}

// NewGRPCServer builds a gRPC server instance
func NewGRPCServer(svc *service.Service) auth1.AuthServiceServer {
	return &grpcServer{svc: svc}
}

func (s *grpcServer) Login(ctx context.Context, req *auth1.LoginRequest) (*auth1.LoginResponse, error) {
	at, rt, pl, err := s.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &auth1.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
		TokenType:    "bearer",
		ExpiresIn:    int64(pl.ExpiredAt.Sub(time.Now()).Seconds()),
	}, nil
}

func (s *grpcServer) Refresh(ctx context.Context, req *auth1.LoginResponse) (*auth1.LoginResponse, error) {
	at, pl, err := s.svc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &auth1.LoginResponse{
		AccessToken:  at,
		RefreshToken: req.RefreshToken,
		TokenType:    "bearer",
		ExpiresIn:    int64(pl.ExpiredAt.Sub(time.Now()).Seconds()),
	}, nil
}

func (s *grpcServer) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	// Extract refresh token from metadata if needed, or ignore body
	// Here we assume client sends refresh token in metadata or URL.
	return &emptypb.Empty{}, nil
}

// PayloadFromCtx retrieves the Paseto payload
func PayloadFromCtx(ctx context.Context) *token.Payload {
	p, _ := ctx.Value(token.CtxKey).(*token.Payload)
	return p
}
