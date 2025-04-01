package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExtractRedisDetails extracts both the address and DB index from a Redis DSN
// Returns address, dbIndex, error
func ExtractRedisDetails(redisDsn string) (string, int, error) {
	// Regular expression to capture the DB index at the end of the DSN (e.g., redis://localhost:6379/1)
	re := regexp.MustCompile(`^(.*)/(\d+)$`)
	matches := re.FindStringSubmatch(redisDsn)

	if len(matches) < 3 {
		return "", 0, fmt.Errorf("invalid Redis DSN: no DB index found")
	}

	address := matches[1]
	// Remove any trailing slash from address
	address = strings.TrimSuffix(address, "/")

	// Convert DB index from string to integer
	dbIndex, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, fmt.Errorf("invalid Redis DB index: %w", err)
	}

	return address, dbIndex, nil
}
