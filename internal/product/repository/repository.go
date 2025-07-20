package repository

import (
	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
)

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id string) (*model.Product, error)
	List(offset, limit int) ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
}
