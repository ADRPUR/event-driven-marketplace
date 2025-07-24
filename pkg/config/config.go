package config

// Tiny helper for loading runtime configuration. It first loads a `.env` file
// (if present) via `godotenv`, then reads environment variables, falling back
// to sensible defaults for local development. In production, just set the envs
// or extend this file with a dedicated configuration library.
// All comments are in English, as required.

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config aggregates all runtime settings for the Product‑Service.
// Extend this struct with additional fields (Kafka, Redis, etc.) as needed.

type Config struct {
	ProdHTTPAddr string // HTTP listen address, default ":8080"
	ProdGRPCAddr string // gRPC listen address, default ":50051"
	AuthHTTPAddr string // HTTP listen address, default ":8090"
	AuthGRPCAddr string // gRPC listen address, default ":50052"
	DBURL        string // Postgres DSN (required)
	SymmetricKey string // 32‑byte key for Paseto (required)
}

// Load loads .env (when present) and returns a Config struct.
// It terminates the program if mandatory variables are missing.
//
//	HTTP_ADDR      → default ":8080"
//	GRPC_ADDR      → default ":50051"
//	DATABASE_URL   → REQUIRED, no default
//	SYMMETRIC_KEY  → REQUIRED, exactly 32 characters (v2‑local)
func Load() Config {
	// Load .env silently; ignore error when file not found.
	_ = godotenv.Load()

	cfg := Config{
		AuthHTTPAddr: getEnv("AUTH_HTTP_ADDR", ":8090"),
		ProdHTTPAddr: getEnv("PROD_HTTP_ADDR", ":8080"),
		AuthGRPCAddr: getEnv("AUTH_GRPC_ADDR", ":50052"),
		ProdGRPCAddr: getEnv("PROD_GRPC_ADDR", ":50051"),
		DBURL:        mustGetEnv("DATABASE_URL"),
		SymmetricKey: mustGetEnv("SYMMETRIC_KEY"),
	}

	if len(cfg.SymmetricKey) != 32 {
		log.Fatalf("SYMMETRIC_KEY must be exactly 32 characters, got %d", len(cfg.SymmetricKey))
	}
	return cfg
}

// getEnv returns the value or a fallback when unset.
func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// mustGetEnv fetches an env var or terminates the program if missing.
func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("missing %s environment variable", key)
	return "" // unreachable
}
