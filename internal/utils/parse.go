package utils

import (
	"strconv"
)

func ParseIntDefault(value string, defaultValue int) int {
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

func ParsePagination(pageStr string, limitStr string) (int, int) {
	page := ParseIntDefault(pageStr, 1)
	limit := ParseIntDefault(limitStr, 20)

	limit = min(20, limit)

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	return limit, offset
}
