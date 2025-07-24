.PHONY: generate migrate test run new_migration

# ──────────────────────────────────────────────────────────────────────────────
# Root Makefile
# ──────────────────────────────────────────────────────────────────────────────

SHELL := /bin/bash

# If .env exists, load it so DATABASE_URL (and any other) get defined
ifneq ("$(wildcard .env)","")
    include .env
    export
endif

GO           := go
PROTO_DIR    := api/proto
TIMESTAMP := $(shell date +%Y%m%d%H%M%S)
NAME      ?= unnamed

# Fail if DATABASE_URL is not set
ifndef DATABASE_URL
$(error "DATABASE_URL is not set. Please define it in .env or your shell.")
endif

# generate Go code from .proto via Buf
generate:
	buf generate . --path $(PROTO_DIR)

# runs DB migrations (example for product‑service)
migrate:
	@echo "🟢 Running migrations..."
	migrate \
	  -path migrations/sql \
	  -database "$(DATABASE_URL)" \
	  up

# run all tests
 test:
	$(GO) test ./... -v -race

# start ProductService (gRPC + REST)
run:
	$(GO) run cmd/product-service/main.go
	$(GO) run cmd/auth-service/main.go

new_migration:
	@mkdir -p migrations/sql
	@touch migrations/sql/$(TIMESTAMP)_$(NAME).up.sql
	@touch migrations/sql/$(TIMESTAMP)_$(NAME).down.sql
	@echo "Created migrations for $(SERVICE): $(TIMESTAMP)_$(NAME).[up|down].sql"