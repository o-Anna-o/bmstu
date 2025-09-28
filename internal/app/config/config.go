package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPass      string
	ServicePort int
}

func NewConfig() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBName:      getEnv("DB_NAME", "DB_RIP"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPass:      getEnv("DB_PASS", "mydbpass"),
		ServicePort: getEnvAsInt("SERVICE_PORT", 8080),
	}

	log.Info("Configuration loaded successfully")
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
