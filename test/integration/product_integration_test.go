package integration

// Integration test for Product repository against a real Postgres instance via testcontainers‑go.

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/repository"
)

// startPostgresContainer boots a Postgres 15 container and returns the handle + DSN.
func startPostgresContainer(ctx context.Context, t *testing.T) (*tcpostgres.PostgresContainer, string) {
	container, err := tcpostgres.RunContainer(ctx,
		tc.WithImage("postgres:15-alpine"),
		tcpostgres.WithDatabase("marketplace"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tc.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/marketplace?sslmode=disable", host, port.Port())
	return container, dsn
}

func TestProductCRUD(t *testing.T) {
	ctx := context.Background()

	container, dsn := startPostgresContainer(ctx, t)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	// connect GORM
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// auto‑migrate schema
	require.NoError(t, db.AutoMigrate(&model.Product{}))

	repo := repository.NewGormRepository(db)

	// ---- CREATE ----
	prod := &model.Product{
		ID:          uuid.New(),
		Name:        "Test product",
		Description: "From integration test",
		Price:       199.99,
	}
	require.NoError(t, repo.Create(ctx, prod))

	// ---- READ ----
	got, err := repo.GetByID(ctx, prod.ID.String())
	require.NoError(t, err)
	require.Equal(t, prod.Name, got.Name)

	// ---- UPDATE ----
	prod.Price = 149.99
	require.NoError(t, repo.Update(ctx, prod))

	// ---- DELETE ----
	require.NoError(t, repo.Delete(ctx, prod.ID.String()))
}
