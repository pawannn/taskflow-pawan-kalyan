package utils

import "time"

func ParseDate(dateStr *string) *time.Time {
	if dateStr == nil {
		return nil
	}

	t, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return nil
	}

	return &t
}
