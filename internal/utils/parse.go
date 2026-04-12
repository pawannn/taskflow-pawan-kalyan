package utils

import (
	"strconv"
)

// ParsePagination calculates limit and offset from page and limit strings with defaults and bounds.
func ParsePagination(pageStr string, limitStr string) (int, int) {
	page := parseIntDefault(pageStr, 1)
	limit := parseIntDefault(limitStr, 20)

	limit = min(20, limit)

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	return limit, offset
}

// ParseIntDefault parses a string to int or returns a default value if invalid or non-positive.
func parseIntDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	if parsed <= 0 {
		return defaultValue
	}

	return parsed
}
