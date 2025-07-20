FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o product-service ./cmd/product-service

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/product-service /product-service
EXPOSE 8080
ENTRYPOINT ["/product-service"]