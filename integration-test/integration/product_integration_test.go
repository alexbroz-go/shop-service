package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "product/proto"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestProductServiceIntegration(t *testing.T) {
	ctx := context.Background()
	netName := "product_net"

	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           netName,
			CheckDuplicate: true,
		},
	})
	require.NoError(t, err)
	defer network.Remove(ctx)

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:15",
			Env: map[string]string{
				"POSTGRES_USER":     "user",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "product",
			},
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
			Networks:     []string{netName},
			NetworkAliases: map[string][]string{
				netName: {"product-db"},
			},
		},
		Started: true,
	})
	require.NoError(t, err)
	defer postgres.Terminate(ctx)

	product, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "product-image:latest",
			Env: map[string]string{
				"DB_HOST":     "product-db",
				"DB_PORT":     "5432",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "product",
			},
			ExposedPorts: []string{"50051/tcp"},
			WaitingFor:   wait.ForListeningPort("50051/tcp"),
			Networks:     []string{netName},
			NetworkAliases: map[string][]string{
				netName: {"product"},
			},
		},
		Started: true,
	})
	require.NoError(t, err)
	defer product.Terminate(ctx)

	host, err := product.Host(ctx)
	require.NoError(t, err)
	port, err := product.MappedPort(ctx, "50051")
	require.NoError(t, err)
	addr := fmt.Sprintf("%s:%s", host, port.Port())

	time.Sleep(2 * time.Second)

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	addResp, err := client.AddProduct(ctx, &pb.AddProductRequest{
		Title:       "Тестовый товар",
		Description: "Описание тестового товара",
		Price:       123.45,
	})
	require.NoError(t, err)
	require.NotZero(t, addResp.ProductId)

	getResp, err := client.GetProduct(ctx, &pb.GetProductRequest{ProductId: addResp.ProductId})
	require.NoError(t, err)
	require.Equal(t, "Тестовый товар", getResp.Title)
	require.Equal(t, "Описание тестового товара", getResp.Description)
	require.Equal(t, 123.45, getResp.Price)
}
