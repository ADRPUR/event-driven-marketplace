package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	}), &gorm.Config{})
	assert.NoError(t, err)
	return db, mock, func() { _ = sqlDB.Close() }
}

func TestGormRepository_CreateUser(t *testing.T) {
	db, mock, cleanup := setupDB(t)
	defer cleanup()
	repo := repository.NewGormRepository(db)

	user := &model.User{ID: uuid.New(), Email: "test@abc.com", PasswordHash: "hash", Role: "user"}
	details := &model.UserDetails{UserID: user.ID, FirstName: "T", Phone: "123"}

	mock.ExpectBegin()
	// users: id, email, password_hash, role, created_at, updated_at, deleted_at
	mock.ExpectExec("INSERT INTO \"users\"").
		WithArgs(user.ID, user.Email, user.PasswordHash, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// user_details: user_id, first_name, last_name, date_of_birth, phone, address, photo_path, thumbnail_path, created_at, updated_at, deleted_at
	mock.ExpectExec("INSERT INTO \"user_details\"").
		WithArgs(
			details.UserID,    // 1
			details.FirstName, // 2
			details.LastName,  // 3
			sqlmock.AnyArg(),  // 4 date_of_birth
			details.Phone,     // 5
			// nothing for $6 (address is NULL)
			details.PhotoPath,     // 7
			details.ThumbnailPath, // 8 ($8)
			sqlmock.AnyArg(),      // 9 created_at ($9)
			sqlmock.AnyArg(),      // 10 updated_at ($10)
			sqlmock.AnyArg(),      // 11 deleted_at ($11)
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user, details)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormRepository_GetByEmail_UserNotFound(t *testing.T) {
	db, mock, cleanup := setupDB(t)
	defer cleanup()
	repo := repository.NewGormRepository(db)

	email := "notfound@abc.com"
	// users: email + limit (GORM automat adaugă LIMIT $2)
	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	u, d, err := repo.GetByEmail(context.Background(), email)
	assert.Nil(t, u)
	assert.Nil(t, d)
	assert.ErrorIs(t, err, repository.ErrUserNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormRepository_GetByEmail_Success(t *testing.T) {
	db, mock, cleanup := setupDB(t)
	defer cleanup()
	repo := repository.NewGormRepository(db)

	id := uuid.New()
	email := "exists@abc.com"

	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(id, email, "hash", "user", time.Now(), time.Now(), nil))
	mock.ExpectQuery("SELECT \\* FROM \"user_details\" WHERE user_id = \\$1 AND \"user_details\"\\.\"deleted_at\" IS NULL ORDER BY \"user_details\"\\.\"user_id\" LIMIT \\$2").
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_id", "first_name", "last_name", "date_of_birth", "phone", "address", "photo_path", "thumbnail_path", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(id, "N", "L", time.Now(), "123", []byte(`{"city":"Bucharest"}`), "/static/photo.jpg", "", time.Now(), time.Now(), nil))

	u, d, err := repo.GetByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, "N", d.FirstName)
	assert.Equal(t, "123", d.Phone)
	assert.NotNil(t, d.Address)
	var addr map[string]any
	_ = json.Unmarshal(d.Address, &addr)
	assert.Equal(t, "Bucharest", addr["city"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormRepository_DeleteUser(t *testing.T) {
	db, mock, cleanup := setupDB(t)
	defer cleanup()
	repo := repository.NewGormRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"user_details\" SET \"deleted_at\"=\\$1 WHERE user_id = \\$2 AND \"user_details\"\\.\"deleted_at\" IS NULL").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE \"users\" SET \"deleted_at\"=\\$1 WHERE id = \\$2 AND \"users\"\\.\"deleted_at\" IS NULL").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
