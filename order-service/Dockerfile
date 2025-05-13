# Компиляция
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o order-service ./cmd/main.go

# Используем минималистичный образ для запуска
FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/order-service .
COPY .env .

CMD ["./order-service"]