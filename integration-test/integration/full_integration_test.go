package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	pbProduct "product/proto"
	pbUser "user-cart-order/proto"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestFullIntegration(t *testing.T) {
	ctx := context.Background()
	networkName := "e2e_test_network"

	// Создаём сеть для контейнеров
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	require.NoError(t, err)
	defer network.Remove(ctx)

	// Запускаем Postgres контейнер
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env:          map[string]string{"POSTGRES_USER": "user", "POSTGRES_PASSWORD": "password", "POSTGRES_DB": "app"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"postgres"},
		},
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: postgresReq, Started: true})
	require.NoError(t, err)
	defer postgres.Terminate(ctx)

	// Запускаем Redis контейнер
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:7.0-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"redis"},
		},
	}
	redis, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: redisReq, Started: true})
	require.NoError(t, err)
	defer redis.Terminate(ctx)

	// Запускаем Product сервис
	productReq := testcontainers.ContainerRequest{
		Image: "product-image:latest",
		Env: map[string]string{
			"DB_HOST":     "postgres",
			"DB_PORT":     "5432",
			"DB_USER":     "user",
			"DB_PASSWORD": "password",
			"DB_NAME":     "app",
		},
		ExposedPorts: []string{"50051/tcp"},
		WaitingFor:   wait.ForListeningPort("50051/tcp"),
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"product"},
		},
	}
	product, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: productReq, Started: true})
	require.NoError(t, err)
	defer product.Terminate(ctx)

	// Запускаем User-Cart-Order сервис
	userReq := testcontainers.ContainerRequest{
		Image: "user-cart-order-image:latest",
		Env: map[string]string{
			"DB_HOST":     "postgres",
			"DB_PORT":     "5432",
			"DB_USER":     "user",
			"DB_PASSWORD": "password",
			"DB_NAME":     "app",
		},
		ExposedPorts: []string{"50052/tcp"},
		WaitingFor:   wait.ForListeningPort("50052/tcp"),
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"user-cart-order"},
		},
	}
	userCartOrder, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: userReq, Started: true})
	require.NoError(t, err)
	defer userCartOrder.Terminate(ctx)

	// Запускаем LLM сервис
	llmReq := testcontainers.ContainerRequest{
		Image: "llm-service-image:latest",
		Env: map[string]string{
			"REDIS_HOST":       "redis",
			"REDIS_PORT":       "6379",
			"YANDEX_API_KEY":   os.Getenv("YANDEX_API_KEY"),
			"YANDEX_FOLDER_ID": os.Getenv("YANDEX_FOLDER_ID"),
			"MODEL_NAME":       "yandexgpt",
		},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog("Uvicorn running on http://0.0.0.0:8000"),
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"llm"},
		},
	}
	llm, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: llmReq, Started: true})
	require.NoError(t, err)
	defer llm.Terminate(ctx)

	// Получаем хост и порт Product сервиса и подключаемся через gRPC
	productHost, err := product.Host(ctx)
	require.NoError(t, err)
	productPort, err := product.MappedPort(ctx, "50051/tcp")
	require.NoError(t, err)
	productAddress := fmt.Sprintf("%s:%s", productHost, productPort.Port())

	productConn, err := grpc.Dial(productAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer productConn.Close()
	productClient := pbProduct.NewProductServiceClient(productConn)

	// Подключаемся к User-Cart-Order сервису
	userHost, err := userCartOrder.Host(ctx)
	require.NoError(t, err)
	userPort, err := userCartOrder.MappedPort(ctx, "50052/tcp")
	require.NoError(t, err)
	userAddress := fmt.Sprintf("%s:%s", userHost, userPort.Port())

	userConn, err := grpc.Dial(userAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer userConn.Close()
	userClient := pbUser.NewUserCartOrderServiceClient(userConn)

	// Регистрируем пользователя
	regResp, err := userClient.Register(ctx, &pbUser.RegisterRequest{
		Login:    "testuser",
		Email:    "test@example.com",
		Password: "secret",
	})
	require.NoError(t, err)
	require.True(t, regResp.Success)

	// Логинимся
	loginResp, err := userClient.Login(ctx, &pbUser.LoginRequest{
		Login:    "testuser",
		Password: "secret",
	})
	require.NoError(t, err)
	userID := loginResp.UserId

	// Добавляем продукт
	addProdResp, err := productClient.AddProduct(ctx, &pbProduct.AddProductRequest{
		Title:       "Шуруповерт",
		Description: "Мощный шуруповерт",
		Price:       99.99,
	})
	require.NoError(t, err)
	prodID := addProdResp.ProductId

	// Добавляем товар в корзину
	_, err = userClient.AddToCart(ctx, &pbUser.AddToCartRequest{
		UserId:    userID,
		ProductId: prodID,
	})
	require.NoError(t, err)

	// Создаём заказ
	orderResp, err := userClient.CreateOrder(ctx, &pbUser.CreateOrderRequest{
		UserId:    userID,
		ProductId: prodID,
		Status:    "created",
	})
	require.NoError(t, err)
	require.NotZero(t, orderResp.Id)

	// Запрос к LLM сервису
	llmHost, err := llm.Host(ctx)
	require.NoError(t, err)
	llmPort, err := llm.MappedPort(ctx, "8000/tcp")
	require.NoError(t, err)
	llmURL := fmt.Sprintf("http://%s:%s/chat", llmHost, llmPort.Port())

	reqBody, err := json.Marshal(map[string]string{
		"user_id":  strconv.Itoa(int(userID)),
		"question": "Расскажи о шуруповерте",
	})
	require.NoError(t, err)

	resp, err := http.Post(llmURL, "application/json", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, 200, resp.StatusCode)
	require.NotEmpty(t, body)
}
