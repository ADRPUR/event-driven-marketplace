package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered account in the marketplace.
// Passwords are stored as bcrypt hashes.
// In production, store hash cost & pepper in config.

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex;size:320;not null"`
	Password  string    `gorm:"size:60;not null"` // bcrypt hash
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
