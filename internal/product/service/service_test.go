package service_test

// Unit‑tests for ProductService using testify/mock.

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/service"
)

// ---- Mock repository ----

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, p *model.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Product), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockRepo) List(ctx context.Context, off, lim int) ([]model.Product, error) {
	args := m.Called(ctx, off, lim)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Product), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockRepo) Update(ctx context.Context, p *model.Product) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockRepo) Delete(ctx context.Context, id string) error { return m.Called(ctx, id).Error(0) }

// ---- Tests ----

func TestCreate_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := service.New(repo)

	repo.On("Create", ctx, mock.AnythingOfType("*model.Product")).Return(nil)

	got, err := svc.Create(ctx, "Laptop", "Gaming laptop", 999.99)

	assert.NoError(t, err)
	assert.Equal(t, "Laptop", got.Name)
	repo.AssertExpectations(t)
}
