# Многоэтапная сборка для оптимизации размера итогового образа

# Этап сборки: используем официальный образ Go для компиляции
FROM golang:1.24-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Устанавливаем protoc и необходимые зависимости
RUN apk add --no-cache protobuf protobuf-dev git

# Устанавливаем protoc плагины для Go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Копируем файлы зависимостей для кэширования
COPY product/go.mod product/go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект
COPY product/ ./

COPY schema-registry/proto/product.proto schema-registry/proto/product.proto
# Компилируем proto-файлы
RUN protoc --go_out=. --go_opt=module=product --go-grpc_out=. --go-grpc_opt=module=product schema-registry/proto/product.proto

# Собираем статический исполняемый файл
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

# Финальный образ: минимальный Alpine
FROM alpine:latest

# Установка рабочей директории
WORKDIR /app

# Копируем собранный бинарник из этапа сборки
COPY --from=builder /app/main .
COPY product/migrations ./migrations

# Открываем порт, на котором будет слушать приложение
EXPOSE 50051

# По умолчанию включаем миграции при старте (для тестов/Testcontainers).
ENV MIGRATE_ON_STARTUP=1
# Команда запуска
CMD ["./main"]

#Команда для генерации прото файлов protoc --proto_path=../schema-registry/proto --go_out=./proto --go-grpc_out=./proto ../schema-registry/proto/product.proto