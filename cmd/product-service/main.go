package main

// Main entry point for the Product‑Service. It boots both an HTTP (Gin) server
// and a gRPC server, sharing the same business‑logic layer.

import (
	"context"
	"errors"
	"github.com/ADRPUR/event-driven-marketplace/internal/middleware"
	grpcHandler "github.com/ADRPUR/event-driven-marketplace/internal/product/handler/grpc"
	"github.com/ADRPUR/event-driven-marketplace/pkg/token"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	productv1 "github.com/ADRPUR/event-driven-marketplace/api/proto/product/v1"

	httphandler "github.com/ADRPUR/event-driven-marketplace/internal/product/handler/http"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/repository"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/service"
	"github.com/ADRPUR/event-driven-marketplace/pkg/config"
	"github.com/ADRPUR/event-driven-marketplace/pkg/database"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// ------------------------------------------------------------------
	// 1. Configuration & DB
	// ------------------------------------------------------------------
	cfg := config.Load() // expects fields: DBURL, HTTPAddr, GRPCAddr

	db, err := database.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// ------------------------------------------------------------------
	// 2. Wiring layers
	// ------------------------------------------------------------------
	repo := repository.NewGormRepository(db)
	svc := service.New(repo) // business‑logic layer

	// Create a Paseto token maker
	maker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		log.Fatalf("paseto maker: %v", err)
	}

	// ------------------------------------------------------------------
	// 3. HTTP server (Gin)
	// ------------------------------------------------------------------
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.AuthMiddleware(maker))

	httphandler.RegisterHTTPRoutes(r, svc)

	if !strings.Contains(cfg.ProdHTTPAddr, ":") {
		cfg.ProdHTTPAddr = ":" + cfg.ProdHTTPAddr
	}

	httpSrv := &http.Server{
		Addr:    cfg.ProdHTTPAddr, // default ":8080"
		Handler: r,
	}

	go func() {
		log.Printf("HTTP server listening on %s", cfg.ProdHTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// ------------------------------------------------------------------
	// 4. gRPC server
	// ------------------------------------------------------------------
	interceptor := grpc.UnaryInterceptor(middleware.AuthUnaryInterceptor(maker))
	grpcSrv := grpc.NewServer(interceptor)
	productv1.RegisterProductServiceServer(grpcSrv, grpcHandler.NewGRPCServer(svc))
	reflection.Register(grpcSrv)

	if !strings.Contains(cfg.ProdGRPCAddr, ":") {
		cfg.ProdGRPCAddr = ":" + cfg.ProdGRPCAddr
	}

	go func() {
		lis, err := net.Listen("tcp", cfg.ProdGRPCAddr) // default ":50051"
		if err != nil {
			log.Fatalf("failed to listen on %s: %v", cfg.ProdGRPCAddr, err)
		}
		log.Printf("gRPC server listening on %s", cfg.ProdGRPCAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// ------------------------------------------------------------------
	// 5. Graceful shutdown
	// ------------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	grpcSrv.GracefulStop()

	log.Println("Product‑Service stopped gracefully")
}
