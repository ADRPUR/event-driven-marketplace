package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
)

var ErrNotFound = errors.New("record not found")

type GormRepo struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepo { return &GormRepo{db: db} }

// Create user
func (r *GormRepo) Create(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetByEmail fetches a user by email
func (r *GormRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

// DeleteUser Delete user (hard-delete)
func (r *GormRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// Save a new refresh token
func (r *GormRepo) Save(ctx context.Context, rt *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

// Get a refresh token record
func (r *GormRepo) Get(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).First(&rt, "token = ?", token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &rt, err
}

// Delete a refresh token
func (r *GormRepo) Delete(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "token = ?", token).Error
}

// DeleteExpired removes old tokens
func (r *GormRepo) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&model.RefreshToken{}).
		Error
}
