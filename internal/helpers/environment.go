package helpers

import (
	"os"
	"strings"
)

var debug = false

func init() {
	if os.Getenv("DJANGOLANG_DEBUG") == "1" {
		debug = true
	}
}

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

func IsDebug() bool {
	return debug
}
