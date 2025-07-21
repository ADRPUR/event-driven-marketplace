.PHONY: generate migrate test run new_migration

GO           := go
PROTO_DIR    := api/proto
TIMESTAMP := $(shell date +%Y%m%d%H%M%S)
SERVICE   ?= product
NAME      ?= unnamed

# generate Go code from .proto via Buf
generate:
	buf generate . --path $(PROTO_DIR)

# runs DB migrations (example for product‑service)
migrate:
	migrate -path migrations/sql -database "$$DATABASE_URL" up

# run all tests
 test:
	$(GO) test ./... -v -race

# start ProductService (gRPC + REST)
run:
	$(GO) run cmd/product-service/main.go

new_migration:
	@mkdir -p migrations/$(SERVICE)/sql
	@touch migrations/$(SERVICE)/sql/$(TIMESTAMP)_$(NAME).up.sql
	@touch migrations/$(SERVICE)/sql/$(TIMESTAMP)_$(NAME).down.sql
	@echo "Created migrations for $(SERVICE): $(TIMESTAMP)_$(NAME).[up|down].sql"