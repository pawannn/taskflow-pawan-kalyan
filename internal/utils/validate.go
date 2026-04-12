package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

var (
	upperRegex   = regexp.MustCompile(`[A-Z]`)
	lowerRegex   = regexp.MustCompile(`[a-z]`)
	numberRegex  = regexp.MustCompile(`[0-9]`)
	specialRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

// ValidateRequired validates required fields and returns field-specific errors.
func ValidateRequired(fields map[string]string) map[string]string {
	errors := map[string]string{}

	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			errors[field] = "is required"
			continue
		}

		if field == "email" {
			if !isValidEmail(value) {
				errors[field] = "is invalid"
			}
		}

		if field == "password" {
			if err := validatePassword(value); err != "" {
				errors[field] = err
			}
		}

	}

	return errors
}

// isValidEmail checks if the given email is valid.
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	if email == "" {
		return false
	}

	_, err := mail.ParseAddress(email)
	return err == nil
}

// validatePassword validates password strength and returns an error message if invalid.
func validatePassword(password string) string {
	if len(password) < 8 {
		return "must be at least 8 characters"
	}

	if len(password) > 64 {
		return "must not exceed 64 characters"
	}

	if !upperRegex.MatchString(password) {
		return "must contain at least one uppercase letter"
	}

	if !lowerRegex.MatchString(password) {
		return "must contain at least one lowercase letter"
	}

	if !numberRegex.MatchString(password) {
		return "must contain at least one number"
	}

	if !specialRegex.MatchString(password) {
		return "must contain at least one special character"
	}

	if regexp.MustCompile(`\s`).MatchString(password) {
		return "must not contain spaces"
	}

	return ""
}
