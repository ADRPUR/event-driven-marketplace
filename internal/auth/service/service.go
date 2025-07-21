package service

// Business‑logic layer for the authentication domain.
// All comments are in English, as required.

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

// Sentinel errors.
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrRefreshNotFound    = errors.New("refresh token not found or expired")
)

// Service orchestrates authentication use‑cases.
// It is transport‑agnostic.
type Service struct {
	users  repository.UserRepository
	tokens repository.TokenRepository
	maker  token.Maker
	atTTL  time.Duration // access‑token TTL
	rtTTL  time.Duration // refresh‑token TTL
}

// New builds a Service.
func New(us repository.UserRepository, ts repository.TokenRepository, mk token.Maker, atTTL, rtTTL time.Duration) *Service {
	return &Service{users: us, tokens: ts, maker: mk, atTTL: atTTL, rtTTL: rtTTL}
}

// Login validates credentials, returns access token, refresh token and payload.
func (s *Service) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, pl *token.Payload, err error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials // hide existence info
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	accessToken, pl, err = s.maker.CreateToken(u.ID, s.atTTL)
	if err != nil {
		return "", "", nil, err
	}

	// generate refresh token (uuid v4 string)
	rtID := uuid.New().String()
	refreshToken = rtID
	if err := s.tokens.Save(ctx, &model.RefreshToken{
		ID:        rtID,
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(s.rtTTL),
	}); err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, pl, nil
}

// Refresh issues a new access token given a valid refresh token.
func (s *Service) Refresh(ctx context.Context, rt string) (string, *token.Payload, error) {
	r, err := s.tokens.Get(ctx, rt)
	if err != nil || r.IsExpired() {
		return "", nil, ErrRefreshNotFound
	}
	return s.mintAccessToken(r.UserID)
}

// Logout deletes the refresh token (single‑device logout).
func (s *Service) Logout(ctx context.Context, rt string) error {
	return s.tokens.Delete(ctx, rt)
}

// helper
func (s *Service) mintAccessToken(userID uuid.UUID) (string, *token.Payload, error) {
	return s.maker.CreateToken(userID, s.atTTL)
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.users.DeleteUser(ctx, id)
}
