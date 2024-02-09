package helpers

import (
	"fmt"
	"os"
	"strings"
)

func GetEnvironmentVariable(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func GetEnvironmentVariableOrDefault(key string, defaultValue string) string {
	value := GetEnvironmentVariable(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func GetDSN() (string, error) {
	postgresUser := strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	if postgresUser == "" {
		postgresUser = "postgres"
	}

	postgresPassword := strings.TrimSpace(os.Getenv("POSTGRES_PASSWORD"))
	if postgresPassword == "" {
		return "", fmt.Errorf("POSTGRES_PASSWORD env var empty or unset")
	}

	postgresHost := strings.TrimSpace(os.Getenv("POSTGRES_HOST"))
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	postgresPort := strings.TrimSpace(os.Getenv("POSTGRES_PORT"))
	if postgresPort == "" {
		postgresPort = "5432"
	}

	postgresDatabase := strings.TrimSpace(os.Getenv("POSTGRES_DB"))
	if postgresDatabase == "" {
		return "", fmt.Errorf("POSTGRES_DB env var empty or unset")
	}

	postgresSSLModeString := "?sslmode=disable"
	if strings.TrimSpace(os.Getenv("POSTGRES_SSLMODE")) == "1" {
		postgresSSLModeString = "?sslmode=enable"
	}

	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v%v",
		postgresUser,
		postgresPassword,
		postgresHost,
		postgresPort,
		postgresDatabase,
		postgresSSLModeString,
	), nil
}

func GetSchema() string {
	postgresSchema := strings.TrimSpace(os.Getenv("POSTGRES_SCHEMA"))
	if postgresSchema == "" {
		postgresSchema = "public"
	}

	return postgresSchema
}
