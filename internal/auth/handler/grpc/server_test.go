package grpc_test

import (
	"context"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"net"
	"testing"
	"time"

	authv1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"
	grpcHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/grpc"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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

// ---- Bufconn helper ----
const bufSize = 1024 * 1024

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) { return lis.Dial() }
}

func startGRPCServer(t *testing.T, svc service.AuthService) (*grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s, grpcHandler.NewGRPCServer(svc))
	go func() {
		_ = s.Serve(lis)
	}()
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	assert.NoError(t, err)
	return conn, func() { s.Stop(); lis.Close() }
}

// ---- TESTS ----

func TestGRPC_Register_OK(t *testing.T) {
	svc := new(mockService)
	svc.On("Register", "gtest@abc.com").Return(nil)

	conn, cleanup := startGRPCServer(t, svc)
	defer cleanup()
	client := authv1.NewAuthServiceClient(conn)

	resp, err := client.Register(context.Background(), &authv1.RegisterRequest{
		Email:    "gtest@abc.com",
		Password: "Abc123!",
		Details:  &authv1.UserDetails{FirstName: "Gina"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGRPC_Login_OK(t *testing.T) {
	svc := new(mockService)
	conn, cleanup := startGRPCServer(t, svc)
	defer cleanup()
	client := authv1.NewAuthServiceClient(conn)

	email := "grpc@abc.com"
	password := "Parola!"
	payload := &token.Payload{UserID: uuid.New(), ExpiredAt: time.Now().Add(time.Hour)}
	svc.On("Login", email, password).Return("at", "rt", "st", payload, nil)
	svc.On("GetUserWithDetails", payload.UserID).Return(&model.User{ID: payload.UserID, Email: email, Role: "user"}, &model.UserDetails{FirstName: "GRpc"}, nil)

	resp, err := client.Login(context.Background(), &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, email, resp.User.Email)
}
