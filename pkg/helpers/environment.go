package helpers

import (
	"log"
	"os"
	"strings"
)

var debug = false

func init() {
	djangolangDebug := GetEnvironmentVariable("DJANGOLANG_DEBUG")
	if djangolangDebug == "" {
		log.Printf("DJANGOLANG_DEBUG empty or unset; defaulted to '0'")
	}

	if djangolangDebug == "1" {
		debug = true
	}
}

func GetEnvironmentVariable(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func GetEnvironmentVariableOrDefault(key string, defaultValue string) string {
	value := GetEnvironmentVariable(key)
	if value == "" {
		log.Printf("%v empty or unset; defaulted to %v", key, defaultValue)
		return defaultValue
	}

	return value
}

func IsDebug() bool {
	return debug
}
