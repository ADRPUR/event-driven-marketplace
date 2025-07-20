package repository

// All comments are in English, as required.
// Package repository defines the persistence layer contract for the Product domain
// and provides an implementation backed by GORM.

import (
	"context"
	"errors"

	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
	"gorm.io/gorm"
)

// ErrNotFound is returned when a product cannot be located in the datastore.
var ErrNotFound = errors.New("product not found")

// Repository abstracts CRUD operations for products.
// Keeping the interface small keeps the service layer decoupled from the
// underlying storage technology (GORM, SQL, etc.).
type Repository interface {
	Create(ctx context.Context, p *model.Product) error
	GetByID(ctx context.Context, id string) (*model.Product, error)
	List(ctx context.Context, offset, limit int) ([]model.Product, error)
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id string) error
}

// gormRepo is a concrete Repository using *gorm.DB.
// It purposefully lives in the same package to stay unexported; callers depend
// on the Repository interface, not the struct.
type gormRepo struct {
	db *gorm.DB
}

// NewGormRepository returns a Repository implemented with GORM.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) Create(ctx context.Context, p *model.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *gormRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {
	var p model.Product
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *gormRepo) List(ctx context.Context, offset, limit int) ([]model.Product, error) {
	var list []model.Product
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormRepo) Update(ctx context.Context, p *model.Product) error {
	tx := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ?", p.ID).
		Updates(map[string]any{
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
