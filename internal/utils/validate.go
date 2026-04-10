package utils

func ValidateRequired(fields map[string]string) map[string]string {
	errors := map[string]string{}

	for field, value := range fields {
		if value == "" {
			errors[field] = "is required"
		}
	}

	return errors
}
