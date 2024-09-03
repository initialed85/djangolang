package helpers

import "fmt"

func GetRedisURL() (string, error) {
	redisURL := GetEnvironmentVariable("REDIS_URL")
	if redisURL == "" {
		return "", fmt.Errorf("REDIS_URL env var empty or unset")
	}

	return redisURL, nil
}
