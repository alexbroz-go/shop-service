# Используем официальный образ Go для сборки
FROM golang:1.24-alpine AS builder

# Установка необходимых инструментов
WORKDIR /app

# Копируем go.mod и go.sum для зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

# Финальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный файл
COPY --from=builder /app/main .

EXPOSE 50052

CMD ["./main"]
