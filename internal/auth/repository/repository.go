package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
)

// UserRepository defines CRUD on users.
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// TokenRepository defines operations on refresh tokens.
type TokenRepository interface {
	Save(ctx context.Context, rt *model.RefreshToken) error
	Get(ctx context.Context, token string) (*model.RefreshToken, error)
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
