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

type GormRepo struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepo { return &GormRepo{db: db} }

// Create ---- UserRepository ----
func (r *GormRepo) Create(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *GormRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *GormRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Hard-delete:
	// return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error

	// Soft-delete (presupunând câmp DeletedAt *time.Time în model.User):
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

// Save ---- TokenRepository ----
func (r *GormRepo) Save(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	rt := model.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(&rt).Error
}

func (r *GormRepo) Delete(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "token = ?", token).Error
}

func (r *GormRepo) Exists(ctx context.Context, token string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("token = ? AND expires_at > NOW()", token).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
