package main

// Entry point for Auth‑Service: starts HTTP (Gin) and gRPC servers.
// Comments in English as required.

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	auth1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"
	grpcHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/grpc"
	httpHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/http"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/repository"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"
	middleware "github.com/ADRPUR/event-driven-marketplace/internal/middleware"
	"github.com/ADRPUR/event-driven-marketplace/pkg/config"
	"github.com/ADRPUR/event-driven-marketplace/pkg/database"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to Postgres
	db, err := database.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// Repos & Service
	repo := repository.NewGormRepository(db)
	svc := service.New(
		repo, repo, // UserRepo, TokenRepo
		token.MustNewPasetoMaker(cfg.SymmetricKey),
		15*time.Minute, 24*time.Hour, // AT TTL, RT TTL
	)

	// ---------- HTTP ----------
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Public routes: /auth/login, /auth/refresh
	httpHandler.RegisterPublicRoutes(r, httpHandler.New(svc))

	// Protected routes: /auth/logout
	authMW := middleware.AuthMiddleware(svc.Maker())
	grp := r.Group("/auth", authMW)
	httpHandler.RegisterProtectedRoutes(grp, httpHandler.New(svc))

	addr := cfg.AuthHTTPAddr
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}
	httpSrv := &http.Server{Addr: addr, Handler: r}

	go func() {
		log.Printf("HTTP listening on %s", addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP error: %v", err)
		}
	}()

	// ---------- gRPC ----------
	grpcAddr := cfg.AuthGRPCAddr
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr = ":" + grpcAddr
	}
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthUnaryInterceptor(svc.Maker())),
	)
	auth1.RegisterAuthServiceServer(grpcSrv, grpcHandler.NewGRPCServer(svc))
	reflection.Register(grpcSrv)

	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("gRPC listen: %v", err)
		}
		log.Printf("gRPC listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("gRPC serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
	log.Println("Auth‑Service stopped gracefully")
}
