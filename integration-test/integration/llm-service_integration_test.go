package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLLMServiceIntegration(t *testing.T) {
	ctx := context.Background()

	// Создаем общую сеть для контейнеров
	networkName := "test-network"
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	require.NoError(t, err)
	defer network.Remove(ctx)

	// Запускаем Redis контейнер
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7.0-alpine",
			ExposedPorts: []string{"6379/tcp"},
			Networks:     []string{networkName},
			NetworkAliases: map[string][]string{
				networkName: {"redis_db"},
			},
			WaitingFor: wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer redisC.Terminate(ctx)

	// Запускаем LLM сервис контейнер
	llmC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "llm-service-image:latest",
			ExposedPorts: []string{"8000/tcp"},
			Env: map[string]string{
				"REDIS_HOST":       "redis_db",
				"REDIS_PORT":       "6379",
				"YANDEX_API_KEY":   os.Getenv("YANDEX_API_KEY"),
				"YANDEX_FOLDER_ID": os.Getenv("YANDEX_FOLDER_ID"),
				"MODEL_NAME":       "yandexgpt",
			},
			Networks:   []string{networkName},
			WaitingFor: wait.ForLog("Uvicorn running on http://0.0.0.0:8000").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer llmC.Terminate(ctx)

	// Получаем хост и порт для доступа к LLM сервису из теста
	host, err := llmC.Host(ctx)
	require.NoError(t, err)
	port, err := llmC.MappedPort(ctx, "8000")
	require.NoError(t, err)

	url := fmt.Sprintf("http://%s:%s/chat", host, port.Port())

	// Формируем тело запроса
	body := `{
        "user_id": "testuser",
        "question": "Расскажи о дрелях."
    }`

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	t.Logf("Response status: %d", resp.StatusCode)
	t.Logf("Response body: %s", string(respBody))

	require.Equal(t, 200, resp.StatusCode, "expected status 200 but got %d", resp.StatusCode)
}
