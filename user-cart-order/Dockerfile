# Используем официальный образ Go для сборки
FROM golang:1.24-alpine AS builder

# Установка необходимых инструментов
WORKDIR /app

# Устанавливаем protoc и необходимые зависимости
RUN apk add --no-cache protobuf protobuf-dev git

# Устанавливаем protoc плагины для Go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Копируем go.mod и go.sum для зависимостей
COPY user-cart-order/go.mod user-cart-order/go.sum ./

# Загружаем зависимости
RUN go mod download

# Создаем структуру директорий для proto-файлов
RUN mkdir -p schema-registry/proto

# Копируем proto-файлы из schema-registry
COPY schema-registry/proto/user-cart-order.proto schema-registry/proto/

# Копируем остальные файлы проекта
COPY user-cart-order/ ./

# Компилируем proto-файлы
RUN protoc --go_out=. --go-grpc_out=. schema-registry/proto/user-cart-order.proto

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

# Финальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный файл
COPY --from=builder /app/main .

EXPOSE 50052

CMD ["./main"]