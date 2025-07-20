.PHONY: generate migrate test run

GO           := go
PROTO_DIR    := api/proto

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