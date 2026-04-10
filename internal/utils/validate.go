package utils

import (
	"net/mail"
	"strings"
)

func ValidateRequired(fields map[string]string) map[string]string {
	errors := map[string]string{}

	for field, value := range fields {
		if field == "email" {
			if !IsValidEmail(value) {
				errors[field] = "is required"
			}

			continue
		}

		if value == "" {
			errors[field] = "is required"
		}
	}

	return errors
}

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	if email == "" {
		return false
	}

	_, err := mail.ParseAddress(email)
	return err == nil
}
