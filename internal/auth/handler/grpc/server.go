package grpc

import (
	"context"
	"encoding/json"
	"gorm.io/datatypes"
	"time"

	auth1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// grpcServer implements auth1.AuthServiceServer
type grpcServer struct {
	auth1.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewGRPCServer(svc service.AuthService) auth1.AuthServiceServer {
	return &grpcServer{svc: svc}
}

// Register ------------------
func (s *grpcServer) Register(ctx context.Context, req *auth1.RegisterRequest) (*auth1.RegisterResponse, error) {
	dob, _ := parseDate(req.Details.DateOfBirth)
	var address datatypes.JSON
	if req.Details.Address != nil {
		addrBytes, err := json.Marshal(req.Details.Address)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid address")
		}
		address = addrBytes
	}
	user := &model.User{
		Email: req.Email,
		Role:  req.Role,
	}
	details := &model.UserDetails{
		FirstName:     req.Details.FirstName,
		LastName:      req.Details.LastName,
		DateOfBirth:   dob,
		Phone:         req.Details.Phone,
		Address:       address,
		PhotoPath:     req.Details.PhotoPath,
		ThumbnailPath: req.Details.ThumbnailPath,
	}
	if err := s.svc.Register(ctx, user, details, req.Password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	return &auth1.RegisterResponse{Id: user.ID.String()}, nil
}

// Login ------------------
func (s *grpcServer) Login(ctx context.Context, req *auth1.LoginRequest) (*auth1.LoginResponse, error) {
	at, rt, st, pl, err := s.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}
	user, details, _ := s.svc.GetUserWithDetails(ctx, pl.UserID)
	return &auth1.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
		SessionToken: st,
		ExpiresAt:    pl.ExpiredAt.Unix(),
		User:         toProtoUser(user, details),
	}, nil
}

// Refresh ------------------
func (s *grpcServer) Refresh(ctx context.Context, req *auth1.RefreshRequest) (*auth1.RefreshResponse, error) {
	at, pl, err := s.svc.Refresh(ctx, req.SessionToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}
	return &auth1.RefreshResponse{
		AccessToken: at,
		ExpiresAt:   pl.ExpiredAt.Unix(),
	}, nil
}

// Logout ------------------
func (s *grpcServer) Logout(ctx context.Context, req *auth1.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.svc.Logout(ctx, req.SessionToken); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

// Me ------------------
func (s *grpcServer) Me(ctx context.Context, _ *auth1.MeRequest) (*auth1.MeResponse, error) {
	payload, ok := ctx.Value(token.CtxKey).(*token.Payload)
	if !ok || payload == nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	user, details, err := s.svc.GetUserWithDetails(ctx, payload.UserID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &auth1.MeResponse{
		User: toProtoUser(user, details),
	}, nil
}

// UpdateUserDetails ------------------
func (s *grpcServer) UpdateUserDetails(ctx context.Context, req *auth1.UpdateUserDetailsRequest) (*emptypb.Empty, error) {
	payload, ok := ctx.Value(token.CtxKey).(*token.Payload)
	if !ok || payload == nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	dob, _ := parseDate(req.Details.DateOfBirth)

	var address datatypes.JSON
	if req.Details.Address != nil {
		addrBytes, err := json.Marshal(req.Details.Address)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid address")
		}
		address = addrBytes
	}
	details := &model.UserDetails{
		UserID:        payload.UserID,
		FirstName:     req.Details.FirstName,
		LastName:      req.Details.LastName,
		DateOfBirth:   dob,
		Phone:         req.Details.Phone,
		Address:       address,
		PhotoPath:     req.Details.PhotoPath,
		ThumbnailPath: req.Details.ThumbnailPath,
	}
	if err := s.svc.UpdateUserDetails(ctx, payload.UserID, details); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

// ChangePassword ------------------
func (s *grpcServer) ChangePassword(ctx context.Context, req *auth1.ChangePasswordRequest) (*emptypb.Empty, error) {
	payload, ok := ctx.Value(token.CtxKey).(*token.Payload)
	if !ok || payload == nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	if err := s.svc.ChangePassword(ctx, payload.UserID, req.OldPassword, req.NewPassword); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "%v", err)
	}
	return &emptypb.Empty{}, nil
}

// UploadPhoto ------------------
func (s *grpcServer) UploadPhoto(ctx context.Context, req *auth1.UploadPhotoRequest) (*auth1.UploadPhotoResponse, error) {
	payload, ok := ctx.Value(token.CtxKey).(*token.Payload)
	if !ok || payload == nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	photoPath, thumbPath, err := s.svc.UploadPhoto(ctx, payload.UserID, req.PhotoData, req.FileExt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &auth1.UploadPhotoResponse{
		PhotoPath:     photoPath,
		ThumbnailPath: thumbPath,
	}, nil
}

// --------- Helpers ---------
func toProtoUser(u *model.User, d *model.UserDetails) *auth1.User {
	if u == nil {
		return nil
	}
	details := &auth1.UserDetails{}
	if d != nil {
		details.FirstName = d.FirstName
		details.LastName = d.LastName
		details.DateOfBirth = d.DateOfBirth.Format("2006-01-02")
		details.Phone = d.Phone
		details.PhotoPath = d.PhotoPath
		details.ThumbnailPath = d.ThumbnailPath

		// Conversion datatypes.JSON → map[string]string
		var addr map[string]string
		if d.Address != nil && len(d.Address) > 0 {
			_ = json.Unmarshal(d.Address, &addr)
		}
		details.Address = addr
	}
	return &auth1.User{
		Id:      u.ID.String(),
		Email:   u.Email,
		Role:    u.Role,
		Details: details,
	}
}

func parseDate(s string) (t time.Time, err error) {
	if s == "" {
		return
	}
	t, err = time.Parse("2006-01-02", s)
	return
}
