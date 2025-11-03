# user-cart-order/Dockerfile
# Многоэтапная сборка для оптимизации размера итогового образа

# Используем официальный образ Go для сборки
FROM golang:1.24-alpine AS builder

# Установка необходимых инструментов
WORKDIR /app

# Устанавливаем protoc и необходимые зависимости
# Это нужно для компиляции proto-файлов в Go код
RUN apk add --no-cache protobuf protobuf-dev git

# Устанавливаем protoc плагины для Go
# Эти плагины необходимы для генерации Go кода из proto-файлов
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="$PATH:$(go env GOPATH)/bin"
# Копируем go.mod и go.sum для зависимостей
# Это позволяет Docker кэшировать слои с зависимостями
# и не скачивать их заново, если не изменились
COPY user-cart-order/go.mod user-cart-order/go.sum ./

# Загружаем зависимости
RUN go mod download

# Создаем структуру директорий для proto-файлов
RUN mkdir -p schema-registry/proto

# Копируем proto-файлы из schema-registry
# Важно: контекст сборки должен быть корневым каталогом проекта,
# иначе эта команда не найдет файлы
COPY schema-registry/proto/user-cart-order.proto schema-registry/proto/

# Копируем остальные файлы проекта
COPY user-cart-order/ ./

# Компилируем proto-файлы с правильным выходным путем
# Добавлены важные флаги --go_opt=module и --go-grpc_opt=module,
# которые указывают protoc генерировать код с правильным импорт-путем
# Это решает проблему, когда исходный код не может найти сгенерированные proto-файлы
RUN protoc --go_out=. --go_opt=module=user-cart-order --go-grpc_out=. --go-grpc_opt=module=user-cart-order schema-registry/proto/user-cart-order.proto

# Собираем приложение
# CGO_ENABLED=0 отключает динамическую компоновку для создания статического исполняемого файла
# GOOS=linux указывает целевую ОС (Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

# Финальный образ для запуска
# Используем минимальный образ alpine для уменьшения размера
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный файл из предыдущего этапа
# Многоэтапная сборка позволяет значительно уменьшить размер итогового образа
COPY --from=builder /app/main .

# Указываем порт, который будет прослушивать приложение
EXPOSE 50052

ENV BOOTSTRAP_SCHEMA=1
# Команда для запуска приложения
CMD ["./main"]