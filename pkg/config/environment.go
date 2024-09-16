package config

import (
	"fmt"
	"os"
	"strings"
)

func GetEnvironmentVariableOrDefault(key string, defaultValue string, notes ...string) string {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		value = defaultValue
	}

	suffix := ""
	if value == defaultValue && defaultValue != "" {
		suffix += "(default) "
	}

	if len(notes) > 0 {
		suffix += fmt.Sprintf("(%s) ", notes[0])
	}

	suffix = strings.TrimSpace(suffix)

	if len(suffix) > 0 {
		suffix = fmt.Sprintf(" %s", suffix)
	}

	if value == defaultValue && defaultValue != "" {
		log.Printf("%v=%#+v%s", key, defaultValue, suffix)
	} else {
		log.Printf("%v=%#+v%s", key, value, suffix)
	}

	return value
}

func GetEnvironmentVariable(key string) string {
	return GetEnvironmentVariableOrDefault(key, "")
}

func GetEnvironmentVariableOrFatal(key string, extras ...string) string {
	env := GetEnvironmentVariable(key)
	if env == "" {
		extra := ""

		if len(extras) > 0 {
			extra = fmt.Sprintf("; %v", extras[0])
		}

		log.Fatalf("%v env var empty or unset%v", key, extra)
	}

	return env
}
