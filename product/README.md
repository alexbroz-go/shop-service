# Проект Название:
Создаем Product - один из микросервисов Shop-service

# Язык программирования:
Golang

# Библиотеки:
github.com/lib/pq v1.10.9
google.golang.org/grpc v1.72.0

# Компиляция Proto-файлов

Проект использует gRPC и Protocol Buffers для коммуникации между сервисами. Ниже приведены инструкции по компиляции proto-файлов.

## Необходимые компоненты

Для компиляции proto-файлов необходимы:
1. Protobuf компилятор (protoc)
2. Go плагины для protoc

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Компиляция Proto-файлов

После установки компонентов выполните:

```
# Из корневой директории Product
protoc --proto_path=../schema-registry/proto --go_out=./proto --go-grpc_out=./proto ../schema-registry/proto/user-cart-order.proto
```

## Обновление Proto-файлов

При внесении изменений в proto-файлы необходимо перегенерировать код:

1. Сначала измените файл в schema-registry/proto/product.proto
2. Затем перекомпилируйте proto-файл, используя одну из команд выше
3. Убедитесь, что новые файлы созданы в директории product/proto/

### Важные заметки:

- Всегда храните исходные proto-файлы в директории schema-registry/proto/
- Сгенерированные Go-файлы хранятся в директории product/proto/
- При изменении API (добавлении новых методов, изменении структуры сообщений) обязательно обновите реализацию сервиса в директории product/service/
- Избегайте изменений, нарушающих обратную совместимость, если это возможно
