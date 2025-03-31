package helpers

import (
	"fmt"
	"regexp"
	"strconv"
)

func ExtractRedisDB(redisDsn string) (int, error) {
	// Regular expression to capture the DB index at the end of the DSN (e.g., redis://localhost:6379/1)
	re := regexp.MustCompile(`/(\d+)$`)
	matches := re.FindStringSubmatch(redisDsn)

	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid Redis DSN: no DB index found")
	}

	// Convert DB index from string to integer
	dbIndex, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid Redis DB index: %w", err)
	}

	return dbIndex, nil
}
