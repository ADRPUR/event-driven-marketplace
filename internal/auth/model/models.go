// internal/auth/model/models.go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User contains authentication-related data.
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Email        string         `gorm:"unique;not null"`
	PasswordHash string         `gorm:"not null"`
	Role         string         `gorm:"default:user;not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// UserDetails stores additional personal user information.
type UserDetails struct {
	UserID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName     string
	LastName      string
	DateOfBirth   time.Time
	Phone         string
	Address       datatypes.JSON
	PhotoPath     string
	ThumbnailPath string
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// Session stores session data for authenticated users.
type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
