package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "user-cart-order/proto"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUserCartOrderIntegration(t *testing.T) {
	ctx := context.Background()

	// --- Создаем отдельную сеть для контейнеров ---
	networkName := "test_network"
	networkReq := testcontainers.NetworkRequest{
		Name:           networkName,
		CheckDuplicate: true,
	}
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: networkReq,
	})
	if err != nil {
		t.Fatal("failed to create network:", err)
	}
	defer network.Remove(ctx)

	// --- Запускаем контейнер с Postgres ---
	postgresReq := testcontainers.ContainerRequest{
		Image: "postgres:15",
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "user_cart_order",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"postgres"},
		},
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		t.Fatal("failed to start postgres container:", err)
	}
	defer postgresC.Terminate(ctx)

	// --- Запускаем контейнер с user-cart-order ---
	userCartOrderReq := testcontainers.ContainerRequest{
		Image: "user-cart-order-image:latest", // твой собранный образ
		Env: map[string]string{
			"DB_HOST":     "postgres", // имя в сети
			"DB_PORT":     "5432",
			"DB_USER":     "user",
			"DB_PASSWORD": "password",
			"DB_NAME":     "user_cart_order",
		},
		ExposedPorts: []string{"50052/tcp"},
		WaitingFor:   wait.ForLog("50052"), // поменяй на лог запуска твоего сервиса
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"user-cart-order"},
		},
	}
	userCartOrderC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: userCartOrderReq,
		Started:          true,
	})
	if err != nil {
		t.Fatal("failed to start user-cart-order container:", err)
	}
	defer userCartOrderC.Terminate(ctx)

	// Получаем хост и порт для подключения к gRPC
	host, err := userCartOrderC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := userCartOrderC.MappedPort(ctx, "50052")
	if err != nil {
		t.Fatal(err)
	}
	grpcAddr := fmt.Sprintf("%s:%s", host, port.Port())

	// Немного подождем для надёжности
	time.Sleep(3 * time.Second)

	// Подключаемся к gRPC сервису
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		t.Fatal("failed to dial grpc:", err)
	}
	defer conn.Close()

	client := pb.NewUserCartOrderServiceClient(conn)

	// 1. Регистрация
	regResp, err := client.Register(ctx, &pb.RegisterRequest{
		Login:    "testuser",
		Email:    "test@example.com",
		Password: "password",
	})
	if err != nil || !regResp.Success {
		t.Fatal("registration failed:", regResp.GetMessage(), err)
	}

	// 2. Логин
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Login:    "testuser",
		Password: "password",
	})
	if err != nil || !loginResp.Success {
		t.Fatal("login failed:", loginResp.GetMessage(), err)
	}
	userID := loginResp.UserId

	// 3. Добавление товара в корзину
	_, err = client.AddToCart(ctx, &pb.AddToCartRequest{
		UserId:    userID,
		ProductId: 1, // допустим, товар с ID 1
	})
	if err != nil {
		t.Fatal("add to cart failed:", err)
	}

	// 4. Создание заказа
	orderResp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId:    userID,
		ProductId: 1,
		Status:    "created",
	})
	if err != nil {
		t.Fatal("create order failed:", err)
	}
	fmt.Println("Order created with ID:", orderResp.Id)

	// 5. Получение истории заказов
	history, err := client.GetOrderHistory(ctx, &pb.UserRequest{
		UserId: userID,
	})
	if err != nil {
		t.Fatal("get order history failed:", err)
	}
	if len(history.Orders) != 1 {
		t.Fatal("order history should contain 1 order")
	}
}
