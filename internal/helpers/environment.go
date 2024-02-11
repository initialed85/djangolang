package helpers

import (
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
