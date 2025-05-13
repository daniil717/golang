# Используем официальный образ Go для сборки
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем бинарный файл
RUN go build -o user-service ./cmd/main.go

# Используем минимальный образ для запуска
FROM debian:bookworm
# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/user-service .

# Копируем файл окружения
COPY .env .

# Указываем порт, который будет использоваться
EXPOSE 50051

# Команда для запуска приложения
CMD ["./user-service"]