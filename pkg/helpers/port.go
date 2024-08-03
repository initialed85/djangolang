package helpers

import (
	"fmt"
	"strconv"
)

func GetPort() (uint16, error) {
	rawPort := GetEnvironmentVariableOrDefault("PORT", "7070")

	port, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("failed to get / parse PORT env var: %v", err)
	}

	return uint16(port), nil
}
