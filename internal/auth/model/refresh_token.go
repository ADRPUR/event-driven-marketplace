package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	Token     string    `gorm:"uniqueIndex;size:128"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// IsExpired returns true if the token has expired.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}
