package service

// All comments are in English, as required.
// Package service holds the business‑logic for the Product domain. It is oblivious
// to the transport (HTTP, gRPC) and to the persistence details.

import (
	"context"
	"errors"

	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/repository"
	"github.com/google/uuid"
)

// ErrNotFound is returned when the requested product does not exist.
var ErrNotFound = errors.New("product not found")

// ProductService is the façade exposed to the transport layers.
// It orchestrates validation and delegates persistence to the repository layer.
type ProductService struct {
	repo repository.Repository
}

// New returns a new ProductService.
func New(repo repository.Repository) *ProductService {
	return &ProductService{repo: repo}
}

// Create inserts a new product and returns the persisted entity.
func (s *ProductService) Create(ctx context.Context, name, description string, price float64) (*model.Product, error) {
	p := &model.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Get retrieves a product by ID.
func (s *ProductService) Get(ctx context.Context, id string) (*model.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

// List returns a paginated set of products.
func (s *ProductService) List(ctx context.Context, page, pageSize int) ([]model.Product, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}

// Update modifies an existing product.
func (s *ProductService) Update(ctx context.Context, id string, name, description *string, price *float64) (*model.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name != nil {
		p.Name = *name
	}
	if description != nil {
		p.Description = *description
	}
	if price != nil {
		p.Price = *price
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Delete removes a product.
func (s *ProductService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
