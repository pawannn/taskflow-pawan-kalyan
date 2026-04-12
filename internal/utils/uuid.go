package utils

import "github.com/google/uuid"

// GenerateUUID returns a new UUID string.
func GenerateUUID() string {
	uuid := uuid.NewString()
	return uuid
}

// IsValidUUID checks whether a string is a valid UUID.
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
