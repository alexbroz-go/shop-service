package config

import (
	"os"
	"sync"
)

// Config holds all application configuration
type Config struct {
	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Add other settings as needed
}

var (
	once     sync.Once
	instance *Config
)

// GetConfig returns a singleton Config instance
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			// Database settings
			DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
			DBPort:     getEnvOrDefault("DB_PORT", "5432"),
			DBUser:     getEnvOrDefault("DB_USER", "postgres"),
			DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:     getEnvOrDefault("DB_NAME", "postgres"),

			// Add other settings here
		}
	})
	return instance
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
