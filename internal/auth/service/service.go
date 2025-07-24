package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/repository"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
)

// Sentinel errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrSessionNotFound    = errors.New("session not found or expired")
	ErrUserNotFound       = errors.New("user not found")
)

// Service implements authentication and user management logic.
type Service struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
	maker    token.Maker
	atTTL    time.Duration // access-token TTL
	rtTTL    time.Duration // refresh-token TTL
}

// New returns a new Service.
func New(users repository.UserRepository, sessions repository.SessionRepository, maker token.Maker, atTTL, rtTTL time.Duration) *Service {
	return &Service{users: users, sessions: sessions, maker: maker, atTTL: atTTL, rtTTL: rtTTL}
}

// Register creates a new user with details and hashes the password.
func (s *Service) Register(ctx context.Context, user *model.User, details *model.UserDetails, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	user.ID = uuid.New()
	details.UserID = user.ID
	return s.users.Create(ctx, user, details)
}

// Login checks credentials, returns tokens and session info.
func (s *Service) Login(ctx context.Context, email, password string) (accessToken, refreshToken, sessionToken string, pl *token.Payload, err error) {
	user, _, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", "", nil, ErrUserNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", "", nil, ErrInvalidCredentials
	}
	accessToken, pl, err = s.maker.CreateToken(user.ID, s.atTTL)
	if err != nil {
		return "", "", "", nil, err
	}
	// Generate session
	session := &model.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(s.rtTTL),
		CreatedAt: time.Now(),
	}
	if err := s.sessions.CreateSession(ctx, session); err != nil {
		return "", "", "", nil, err
	}
	refreshToken = session.Token
	sessionToken = session.Token
	return accessToken, refreshToken, sessionToken, pl, nil
}

// Refresh checks session and issues new access token
func (s *Service) Refresh(ctx context.Context, sessionToken string) (string, *token.Payload, error) {
	session, err := s.sessions.GetSessionByToken(ctx, sessionToken)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		return "", nil, ErrSessionNotFound
	}
	return s.maker.CreateToken(session.UserID, s.atTTL)
}

// Logout deletes session
func (s *Service) Logout(ctx context.Context, sessionToken string) error {
	return s.sessions.DeleteSession(ctx, sessionToken)
}

// GetUserWithDetails retrieves full user profile by id
func (s *Service) GetUserWithDetails(ctx context.Context, id uuid.UUID) (*model.User, *model.UserDetails, error) {
	return s.users.GetByID(ctx, id)
}

// UpdateUserDetails updates the user's detailed data.
func (s *Service) UpdateUserDetails(ctx context.Context, userID uuid.UUID, details *model.UserDetails) error {
	return s.users.Update(ctx, &model.User{ID: userID}, details)
}

// ChangePassword changes the user's password after verifying the old one.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, _, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)) != nil {
		return ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return s.users.Update(ctx, user, nil)
}

// UploadPhoto save the file and update the path in UserDetails.
func (s *Service) UploadPhoto(ctx context.Context, userID uuid.UUID, data []byte, ext string) (string, string, error) {
	photoDir := "./static/photos/users"
	if err := os.MkdirAll(photoDir, 0755); err != nil {
		return "", "", err
	}
	filename := userID.String() + ext
	savePath := filepath.Join(photoDir, filename)
	if err := os.WriteFile(savePath, data, 0644); err != nil {
		return "", "", err
	}
	// TODO: thumbnail generation (for demo, leave empty)
	photoPath := "/static/photos/users" + filename
	thumbPath := ""
	_, details, err := s.GetUserWithDetails(ctx, userID)
	if err == nil && details != nil {
		details.PhotoPath = photoPath
		details.ThumbnailPath = thumbPath
		_ = s.UpdateUserDetails(ctx, userID, details)
	}
	return photoPath, thumbPath, nil
}
