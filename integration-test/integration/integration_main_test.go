package integration_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("./.env")
	if err != nil {
		panic("Ошибка загрузки .env файла: " + err.Error())
	}
	os.Exit(m.Run())
}
