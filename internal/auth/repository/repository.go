package repository

import (
	"context"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/google/uuid"
)

// UserRepository manages users and user_details.
type UserRepository interface {
	Create(ctx context.Context, user *model.User, details *model.UserDetails) error
	GetByEmail(ctx context.Context, email string) (*model.User, *model.UserDetails, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, *model.UserDetails, error)
	Update(ctx context.Context, user *model.User, details *model.UserDetails) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SessionRepository manages user sessions.
type SessionRepository interface {
	CreateSession(ctx context.Context, s *model.Session) error
	GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	DeleteSession(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	CleanupExpired(ctx context.Context) error
}
