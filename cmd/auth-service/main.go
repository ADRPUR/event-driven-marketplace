package main

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
	"github.com/ADRPUR/event-driven-marketplace/internal/middleware"
	"github.com/ADRPUR/event-driven-marketplace/pkg/config"
	"github.com/ADRPUR/event-driven-marketplace/pkg/database"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1) Load config (.env or env variables )
	cfg := config.Load()

	// 2) Connecting to Postgres
	db, err := database.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}

	// 3) Initialization Paseto token maker
	maker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		log.Fatalf("Paseto maker: %v", err)
	}

	// 4) Repository & Service
	repo := repository.NewGormRepository(db)
	svc := service.New(repo, repo, maker, 15*time.Minute, 24*time.Hour)

	// 5) Gin HTTP server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// Expose the photo folder
	r.Static("/static", "./static")

	// Public routes: register, login, refresh
	httpHandler.RegisterPublicRoutes(r, httpHandler.New(svc))
	// Protected routes: /me, /logout, /me/photo etc.
	authMW := middleware.AuthMiddleware(maker)
	protected := r.Group("/", authMW)
	httpHandler.RegisterProtectedRoutes(protected, httpHandler.New(svc))

	httpAddr := cfg.AuthHTTPAddr
	if !strings.Contains(httpAddr, ":") {
		httpAddr = ":" + httpAddr
	}
	httpSrv := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	go func() {
		log.Printf("HTTP server listening on %s", httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 6) gRPC server
	grpcAddr := cfg.AuthGRPCAddr
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr = ":" + grpcAddr
	}
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthUnaryInterceptor(maker)),
	)
	auth1.RegisterAuthServiceServer(grpcSrv, grpcHandler.NewGRPCServer(svc))
	reflection.Register(grpcSrv)

	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
		}
		log.Printf("gRPC server listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// 7) Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	grpcSrv.GracefulStop()
	log.Println("Auth-Service stopped gracefully")
}
