package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	ServerUrl         string
	MetricsReadPeriod int
	ErrorThreshold    int
}

func parseInt(val string) (int, error) {
	return strconv.Atoi(val)
}

func GetConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("No .env file found: %v", err)
	}

	return Config{
		envOrDefault("SERVER_URL", "ws://localhost:8080"),
		envOrDefaultParsed("METRICS_PERIOD", 1000, parseInt),
		envOrDefaultParsed("ERRORS_MAX", 5, parseInt),
	}
}

func envOrDefault(envName string, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	return value
}

func envOrDefaultParsed[T any](envName string, defaultValue T, parseFunc func(string) (T, error)) T {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	parsed, err := parseFunc(value)
	if err != nil {
		fmt.Printf("Failed parsing env %v, value is %v, error: %v", envName, value, err)
		return defaultValue
	}
	return parsed
}
