package utils

import "github.com/google/uuid"

func GenerateUUID() string {
	uuid := uuid.NewString()
	return uuid
}
