package integration

import (
	"context"
	"net"
	"testing"
	"time"

	authv1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"
	grpcHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/grpc"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/repository"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) { return lis.Dial() }
}

func startIntegrationGRPCServer(t *testing.T) (*grpc.ClientConn, func(), service.AuthService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.UserDetails{}, &model.Session{}))
	userRepo := repository.NewGormRepository(db)
	sessionRepo := repository.NewGormRepository(db)
	maker, _ := token.NewPasetoMaker("12345678901234567890123456789012")
	authSvc := service.New(userRepo, sessionRepo, maker, time.Minute, 2*time.Minute)

	// Bufconn + server
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s, grpcHandler.NewGRPCServer(authSvc))
	go func() { _ = s.Serve(lis) }()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithInsecure())
	require.NoError(t, err)
	return conn, func() { s.Stop(); lis.Close() }, authSvc
}

func TestAuthGRPC_Register_Login(t *testing.T) {
	conn, cleanup, _ := startIntegrationGRPCServer(t)
	defer cleanup()
	client := authv1.NewAuthServiceClient(conn)

	// 1. Register
	registerResp, err := client.Register(context.Background(), &authv1.RegisterRequest{
		Email:    "integration@abc.com",
		Password: "Parola123!",
		Details:  &authv1.UserDetails{FirstName: "Inte", LastName: "Test"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, registerResp.Id)

	// 2. Login
	loginResp, err := client.Login(context.Background(), &authv1.LoginRequest{
		Email:    "integration@abc.com",
		Password: "Parola123!",
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.AccessToken)
	require.Equal(t, "integration@abc.com", loginResp.User.Email)
	require.Equal(t, "Inte", loginResp.User.Details.FirstName)
}
