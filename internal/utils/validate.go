package utils

import (
	"net/mail"
	"strings"
)

func ValidateRequired(fields map[string]string) map[string]string {
	errors := map[string]string{}

	for field, value := range fields {
		if value == "" {
			errors[field] = "is required"
			continue
		}

		if field == "email" {
			if !IsValidEmail(value) {
				errors[field] = "is invalid"
			}
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
