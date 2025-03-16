package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) error {
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make([]string, 0)

	for _, e := range validationErrors {
		var message string
		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", e.Field())
		case "uuid":
			message = fmt.Sprintf("%s must be a valid UUID", e.Field())
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
		default:
			message = fmt.Sprintf("%s is invalid", e.Field())
		}
		errorMessages = append(errorMessages, message)
	}

	// Join all error messages into a single string
	return fmt.Errorf("%s", strings.Join(errorMessages, "; "))
}

type ErrorResponse struct {
	Error string `json:"error"`
}
