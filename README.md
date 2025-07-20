# Product Service – Event‑Driven Marketplace

A **Go microservice** that manages product CRUD operations for the Event‑Driven Marketplace. It exposes both **REST (Gin)** and **gRPC** APIs, uses **PostgreSQL (GORM)** for persistence, and secures requests with **Paseto v2‑local tokens**. The service is designed with Clean Architecture principles and comes with unit/integration tests, migrations, and a ready‑to‑use Makefile.

---

## ✨ Key Features

| Layer             | Details                                                       |
| ----------------- | ------------------------------------------------------------- |
| **Transport**     | REST (Gin) + gRPC (reflection enabled)                        |
| **Auth**          | Paseto v2‑local middleware for HTTP & gRPC                    |
| **Storage**       | PostgreSQL 15+ with GORM + SQL migrations                     |
| **Tests**         | Unit (Testify + go‑sqlmock) & Integration (testcontainers‑go) |
| **Observability** | pprof, OpenTelemetry hooks (placeholders)                     |

---

## 🗂️ Project Structure (excerpt)

```
cmd/
  product-service/       # main.go – entry point
internal/
  product/
    model/               # domain models
    repository/          # GORM repository + interface
    service/             # business logic
    handler/
      grpc_server.go     # gRPC transport
      http/handler.go    # Gin transport
pkg/
  config/                # env loader
  token/                 # Paseto maker
  database/              # DB helper
api/proto/…              # Protobuf definitions (product/v1)
migrations/sql/          # Up/Down SQL scripts
Makefile                 # generate, migrate, test, run
```

---

## 🚀 Quick Start

### Prerequisites

* Go ≥ 1.22
* PostgreSQL 15+
* `protoc`, `buf` and the plugins (`protoc-gen-go`, `protoc-gen-go-grpc`)

### 1. Clone & configure

```bash
git clone https://github.com/your‑fork/event-driven-marketplace.git
cd event-driven-marketplace
cp .env.example .env         # or export env vars manually
```

Minimal required env vars:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/marketplace?sslmode=disable
SYMMETRIC_KEY=12345678901234567890123456789012 # 32‑byte key
```

### 2. Generate code & run migrations

```bash
make generate   # buf generate – protobuf → Go
make migrate    # golang‑migrate up
```

### 3. Run the service

```bash
make run        # shorthand for `go run cmd/product-service/main.go`
```

* REST available at **`http://localhost:8080/products`**
* gRPC available at **`localhost:50051`** (use `grpcurl` for quick calls)

---

## 🔐 Authentication

This service expects a **Bearer Paseto** token on each request.

```text
Authorization: Bearer <token>
```

Generate a token in Go:

```go
maker, _ := token.NewPasetoMaker(os.Getenv("SYMMETRIC_KEY"))
userID := uuid.New()
tok, _ := maker.CreateToken(userID, time.Hour)
```

For local testing you can bypass auth by commenting the middleware lines in *main.go*.

---

## 📚 API Overview

### REST (JSON)

| Method | Path            | Description                        |
| ------ | --------------- | ---------------------------------- |
| POST   | `/products`     | Create product                     |
| GET    | `/products/:id` | Get product by ID                  |
| GET    | `/products`     | List products (`?page=&pageSize=`) |
| PUT    | `/products/:id` | Update product                     |
| DELETE | `/products/:id` | Delete product                     |

### gRPC

*Service:* `product.v1.ProductService`

RPCs: `CreateProduct`, `GetProduct`, `ListProducts`, `UpdateProduct`, `DeleteProduct`.

---

## 🧪 Testing

```bash
make test         # runs unit + integration tests
```

Integration tests spin up a temporary Postgres container via `testcontainers‑go` – Docker is required.

---

## 🛠️ Useful Make Targets

| Command         | Action                          |
| --------------- | ------------------------------- |
| `make generate` | Regenerate protobuf & gRPC code |
| `make migrate`  | Apply all SQL migrations up     |
| `make test`     | Run unit & integration tests    |
| `make run`      | Run product‑service locally     |

---

## 🤝 Contributing

1. Fork the repo & create feature branch (`git checkout -b feat/my-feature`)
2. Commit your changes (`git commit -m 'feat: add awesome stuff'`)
3. Push to the branch (`git push origin feat/my-feature`)
4. Open a Pull Request

Please make sure all tests pass and `go vet ./...` returns no issues before submitting.

---

## 📝 License

This project is licensed under the **MIT License**. See `LICENSE` for details.
