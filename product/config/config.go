package config

import (
	"os"
	"sync"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var (
	once     sync.Once
	instance *Config
)

// getConfig возвращает экземпляр конфигурации singleton
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			// Database settings
			DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
			DBPort:     getEnvOrDefault("DB_PORT", "5434"),
			DBUser:     getEnvOrDefault("DB_USER", "user"),
			DBPassword: getEnvOrDefault("DB_PASSWORD", "password"),
			DBName:     getEnvOrDefault("DB_NAME", "product"),
		}
	})
	return instance
}

// Вспомогательная функция для получения переменной окружения со значением по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
