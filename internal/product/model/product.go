package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a marketplace item.
// `gorm` tags map struct fields to table columns.
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// BeforeCreate generates a UUID before inserting the record.
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
