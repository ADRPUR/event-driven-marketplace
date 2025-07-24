package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/repository"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrRefreshNotFound    = errors.New("refresh token not found or expired")
)

type Service struct {
	users  repository.UserRepository
	tokens repository.TokenRepository
	maker  token.Maker
	atTTL  time.Duration
	rtTTL  time.Duration
}

// New constructs the auth service with TTLs.
func New(
	ur repository.UserRepository,
	tr repository.TokenRepository,
	mk token.Maker,
	atTTL, rtTTL time.Duration,
) *Service {
	return &Service{users: ur, tokens: tr, maker: mk, atTTL: atTTL, rtTTL: rtTTL}
}

// Maker exposes the token.Maker for middleware.
func (s *Service) Maker() token.Maker { return s.maker }

// Login returns new access + refresh tokens, and payload.
func (s *Service) Login(
	ctx context.Context,
	email, password string,
) (accessToken string, refreshToken string, pl *token.Payload, err error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	accessToken, pl, err = s.maker.CreateToken(u.ID, s.atTTL)
	if err != nil {
		return
	}

	// generate and store refresh token
	rt := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(s.rtTTL),
	}
	if err = s.tokens.Save(ctx, rt); err != nil {
		return
	}
	return accessToken, rt.Token, pl, nil
}

// Refresh issues a new access token for a valid refresh token.
func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (newAccessToken string, pl *token.Payload, err error) {
	rt, err := s.tokens.Get(ctx, refreshToken)
	if err != nil || rt.IsExpired() {
		return "", nil, ErrRefreshNotFound
	}
	return s.maker.CreateToken(rt.UserID, s.atTTL)
}

// Logout removes a refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.tokens.Delete(ctx, refreshToken)
}
