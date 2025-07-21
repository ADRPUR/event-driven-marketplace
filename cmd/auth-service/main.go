package main

// Entry point for the Auth‑Service. Boots both an HTTP (Gin) server and a gRPC
// server, sharing the same AuthService business layer.
// All comments are in English, as required.

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv1 "github.com/ADRPUR/event-driven-marketplace/api/proto/auth/v1"

	grpcHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/grpc"
	httpHandler "github.com/ADRPUR/event-driven-marketplace/internal/auth/handler/http"

	"github.com/ADRPUR/event-driven-marketplace/internal/auth/repository"
	"github.com/ADRPUR/event-driven-marketplace/internal/auth/service"

	"github.com/ADRPUR/event-driven-marketplace/pkg/config"
	"github.com/ADRPUR/event-driven-marketplace/pkg/database"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1) Load configuration (.env or environment variables)
	cfg := config.Load()

	// 2) Connect to Postgres
	db, err := database.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// 3) Wire repository → service
	repo := repository.NewGormRepository(db)
	svc := service.New(repo)

	// 4) Paseto maker
	maker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		log.Fatalf("paseto maker: %v", err)
	}

	// ------------------------------------------------------------------
	// HTTP server (Gin)
	// ------------------------------------------------------------------
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Public routes (e.g., /login /refresh)
	httpHandler.RegisterPublicRoutes(r, svc, maker)

	// Protected routes example (e.g., /logout)
	authMW := httpHandler.AuthMiddleware(maker)
	protected := r.Group("/", authMW)
	httpHandler.RegisterProtectedRoutes(protected, svc)

	httpSrv := &http.Server{
		Addr:    cfg.HTTPAddr, // default ":8090"
		Handler: r,
	}

	go func() {
		log.Printf("Auth‑HTTP listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listen: %v", err)
		}
	}()

	// ------------------------------------------------------------------
	// gRPC server
	// ------------------------------------------------------------------
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(httpHandler.GRPCAuthInterceptor(maker)),
	)
	authv1.RegisterAuthServiceServer(grpcSrv, grpcHandler.NewGRPCServer(svc, maker))
	reflection.Register(grpcSrv)

	go func() {
		lis, err := net.Listen("tcp", cfg.GRPCAddr) // default ":50052"
		if err != nil {
			log.Fatalf("grpc listen: %v", err)
		}
		log.Printf("Auth‑gRPC listening on %s", cfg.GRPCAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	// ------------------------------------------------------------------
	// Graceful shutdown
	// ------------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	grpcSrv.GracefulStop()
	log.Println("Auth‑Service stopped gracefully")
}
