package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/http"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

// --- Mock service ---
type mockService struct{ mock.Mock }

func (m *mockService) Register(ctx context.Context, user *model.User, details *model.UserDetails, password string) error {
	args := m.Called(user.Email)
	return args.Error(0)
}
func (m *mockService) Login(ctx context.Context, email, password string) (string, string, string, *token.Payload, error) {
	args := m.Called(email, password)
	return args.String(0), args.String(1), args.String(2), args.Get(3).(*token.Payload), args.Error(4)
}
func (m *mockService) Refresh(ctx context.Context, sessionToken string) (string, *token.Payload, error) {
	args := m.Called(sessionToken)
	return args.String(0), args.Get(1).(*token.Payload), args.Error(2)
}
func (m *mockService) Logout(ctx context.Context, sessionToken string) error {
	args := m.Called(sessionToken)
	return args.Error(0)
}
func (m *mockService) GetUserWithDetails(ctx context.Context, id uuid.UUID) (*model.User, *model.UserDetails, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Get(1).(*model.UserDetails), args.Error(2)
}
func (m *mockService) UpdateUserDetails(ctx context.Context, userID uuid.UUID, details *model.UserDetails) error {
	args := m.Called(userID)
	return args.Error(0)
}
func (m *mockService) ChangePassword(ctx context.Context, userID uuid.UUID, old, new string) error {
	args := m.Called(userID, old, new)
	return args.Error(0)
}
func (m *mockService) UploadPhoto(ctx context.Context, userID uuid.UUID, data []byte, ext string) (string, string, error) {
	args := m.Called(userID)
	return args.String(0), args.String(1), args.Error(2)
}

func setupRouter(svc *mockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := httpHandler.New(svc)
	httpHandler.RegisterPublicRoutes(r, h)
	return r
}

func TestRegister_OK(t *testing.T) {
	svc := new(mockService)
	r := setupRouter(svc)
	payload := map[string]any{
		"email":     "test@abc.com",
		"password":  "Abc123!",
		"firstName": "John",
	}
	svc.On("Register", "test@abc.com").Return(nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin_OK(t *testing.T) {
	svc := new(mockService)
	r := setupRouter(svc)
	email := "test@abc.com"
	password := "Abc123!"
	payload := &token.Payload{UserID: uuid.New(), ExpiredAt: time.Now().Add(time.Hour)}
	user := &model.User{ID: payload.UserID, Email: email, Role: "user"}
	details := &model.UserDetails{FirstName: "John", Address: datatypes.JSON([]byte(`{"city":"Chisinau"}`))}
	svc.On("Login", email, password).Return("at", "rt", "st", payload, nil)
	svc.On("GetUserWithDetails", payload.UserID).Return(user, details, nil)

	reqBody := map[string]any{"email": email, "password": password}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, email, resp["user"].(map[string]any)["email"])
}
