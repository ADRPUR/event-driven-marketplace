package repository

import (
	"context"
	_ "database/sql"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	dial := postgres.New(postgres.Config{Conn: sqlDB, DriverName: "postgres"})
	gdb, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	return gdb, mock, func() {
		err := sqlDB.Close()
		if err != nil {
			return
		}
	}
}

func TestCreate_OK(t *testing.T) {
	db, mock, closeFn := setupDB(t)
	defer closeFn()

	repo := NewGormRepository(db)

	p := &model.Product{
		ID:    uuid.New(),
		Name:  "Mouse",
		Price: 25.5,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "products"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), p)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
