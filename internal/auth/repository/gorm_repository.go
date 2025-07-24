package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")

// GormRepository implements UserRepository and SessionRepository
type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// --------------------- UserRepository ----------------------

func (r *GormRepository) Create(ctx context.Context, user *model.User, details *model.UserDetails) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		details.UserID = user.ID
		return tx.Create(details).Error
	})
}

func (r *GormRepository) GetByEmail(ctx context.Context, email string) (*model.User, *model.UserDetails, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, err
	}
	var details model.UserDetails
	if err := r.db.WithContext(ctx).Where("user_id = ?", user.ID).First(&details).Error; err != nil {
		return &user, nil, nil // details can be nil for minimal auth
	}
	return &user, &details, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, *model.UserDetails, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, err
	}
	var details model.UserDetails
	if err := r.db.WithContext(ctx).Where("user_id = ?", user.ID).First(&details).Error; err != nil {
		return &user, nil, nil
	}
	return &user, &details, nil
}

func (r *GormRepository) Update(ctx context.Context, user *model.User, details *model.UserDetails) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
			return err
		}
		if details != nil {
			if err := tx.Model(details).Where("user_id = ?", user.ID).Updates(details).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.UserDetails{}, "user_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.User{}, "id = ?", id).Error
	})
}

// --------------------- SessionRepository ----------------------

func (r *GormRepository) CreateSession(ctx context.Context, s *model.Session) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *GormRepository) DeleteSession(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&model.Session{}).Error
}

func (r *GormRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	var s model.Session
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &s, nil
}

// DeleteByUserID delete all sessions for a user.
func (r *GormRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.Session{}).Error
}

// CleanupExpired delete expired sessions.
func (r *GormRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}

// (Optionally, you can add List, Update etc for sessions and users)
