package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCKS ---

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, user *model.User, details *model.UserDetails) error {
	return m.Called(ctx, user, details).Error(0)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, *model.UserDetails, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*model.User), args.Get(1).(*model.UserDetails), args.Error(2)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, *model.UserDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Get(1).(*model.UserDetails), args.Error(2)
}
func (m *mockUserRepo) Update(ctx context.Context, user *model.User, details *model.UserDetails) error {
	return m.Called(ctx, user, details).Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockSessionRepo struct{ mock.Mock }

func (m *mockSessionRepo) CreateSession(ctx context.Context, s *model.Session) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockSessionRepo) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*model.Session), args.Error(1)
}
func (m *mockSessionRepo) DeleteSession(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockSessionRepo) DeleteByUserID(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockSessionRepo) CleanupExpired(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

type mockTokenMaker struct{ token.Maker }

func (m *mockTokenMaker) CreateToken(userID uuid.UUID, ttl time.Duration) (string, *token.Payload, error) {
	return "at", &token.Payload{UserID: userID, ExpiredAt: time.Now().Add(ttl)}, nil
}

// --- TESTS ---

func TestService_Register_Login_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	user := &model.User{Email: "a@b.com"}
	details := &model.UserDetails{FirstName: "A"}
	userRepo.On("Create", ctx, mock.Anything, mock.Anything).Return(nil)

	assert.NoError(t, svc.Register(ctx, user, details, "Secret123!"))

	hashed, _ := service.HashPassword("Secret123!")
	user.ID = uuid.New()
	user.PasswordHash = hashed
	userRepo.On("GetByEmail", ctx, user.Email).Return(user, details, nil)
	sessionRepo.On("CreateSession", ctx, mock.Anything).Return(nil)

	at, rt, st, payload, err := svc.Login(ctx, user.Email, "Secret123!")
	assert.NoError(t, err)
	assert.NotEmpty(t, at)
	assert.NotEmpty(t, rt)
	assert.NotEmpty(t, st)
	assert.Equal(t, user.ID, payload.UserID)
}

func TestService_Login_WrongPassword(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	user := &model.User{Email: "a@b.com"}
	details := &model.UserDetails{}
	hashed, _ := service.HashPassword("RealPassword!")
	user.ID = uuid.New()
	user.PasswordHash = hashed
	userRepo.On("GetByEmail", ctx, user.Email).Return(user, details, nil)

	_, _, _, _, err := svc.Login(ctx, user.Email, "WrongPassword!")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestService_Refresh_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	userID := uuid.New()
	session := &model.Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "token123",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	sessionRepo.On("GetSessionByToken", ctx, "token123").Return(session, nil)

	at, payload, err := svc.Refresh(ctx, "token123")
	assert.NoError(t, err)
	assert.NotEmpty(t, at)
	assert.Equal(t, userID, payload.UserID)
}

func TestService_Logout_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	sessionRepo.On("DeleteSession", ctx, "tok").Return(nil)

	assert.NoError(t, svc.Logout(ctx, "tok"))
}

func TestService_UpdateUserDetails(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	details := &model.UserDetails{UserID: uuid.New(), FirstName: "Nume"}
	user := &model.User{ID: details.UserID}
	userRepo.On("Update", ctx, user, details).Return(nil)

	assert.NoError(t, svc.UpdateUserDetails(ctx, details.UserID, details))
}

func TestService_ChangePassword_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	userID := uuid.New()
	hashed, _ := service.HashPassword("oldpass")
	user := &model.User{ID: userID, PasswordHash: hashed}
	userRepo.On("GetByID", ctx, userID).Return(user, &model.UserDetails{}, nil)
	userRepo.On("Update", ctx, user, (*model.UserDetails)(nil)).Return(nil)

	assert.NoError(t, svc.ChangePassword(ctx, userID, "oldpass", "newpass"))
}

func TestService_UploadPhoto(t *testing.T) {
	userRepo := new(mockUserRepo)
	sessionRepo := new(mockSessionRepo)
	tokenMaker := new(mockTokenMaker)
	svc := service.New(userRepo, sessionRepo, tokenMaker, time.Minute, time.Hour)

	ctx := context.Background()
	userID := uuid.New()
	userRepo.On("GetByID", ctx, userID).Return(&model.User{ID: userID}, &model.UserDetails{}, nil)
	userRepo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil)

	photo := []byte{1, 2, 3}
	photoPath, thumbPath, err := svc.UploadPhoto(ctx, userID, photo, ".jpg")
	assert.NoError(t, err)
	assert.Contains(t, photoPath, ".jpg")
	assert.Equal(t, "", thumbPath)

	// Cleanup: delete the static folder created by the test
	t.Cleanup(func() {
		_ = os.RemoveAll("./static")
	})
}
