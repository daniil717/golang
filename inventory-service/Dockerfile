# Используем официальный образ Go для сборки
FROM golang:1.23 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем бинарный файл
RUN go build -o inventory-service ./cmd/main.go

# Используем минимальный образ для запуска
FROM debian:bookworm

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/inventory-service .

# Копируем файл окружения
COPY .env .

# Указываем порт, который будет использоваться
EXPOSE 50053

# Команда для запуска приложения
CMD ["./inventory-service"]