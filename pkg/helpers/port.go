package helpers

import (
	"fmt"
	"strconv"

	_helpers "github.com/initialed85/djangolang/internal/helpers"
)

func GetPort() (uint16, error) {
	rawPort := _helpers.GetEnvironmentVariable("PORT")
	port, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("failed to get / parse PORT env var: %v", err)
	}

	return uint16(port), nil
}
