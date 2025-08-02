package application_error

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"justcallmesu.com/rest-api/internal/api/response"
)

func FormatValidationError(passedError error) []response.APIValidationError {
	var validationErrors validator.ValidationErrors

	if errors.As(passedError, &validationErrors) {

		var formattedErrors = make([]response.APIValidationError, len(validationErrors))

		for index, err := range validationErrors {
			formattedErrors[index] = response.APIValidationError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
			}
		}
		return formattedErrors
	}

	return nil
}

func getErrorMessage(err validator.FieldError) string {

	switch err.Tag() {
	case "required":
		return "Field " + err.Field() + " is required"
	case "email":
		return "Field " + err.Field() + " must be a valid email address"
	case "min":
		return "Field " + err.Field() + " must be at least " + err.Param() + " characters long"
	case "max":
		return "Field " + err.Field() + " must be at most " + err.Param() + " characters long"
	case "len":
		return "Field " + err.Field() + " must be exactly " + err.Param() + " characters long"
	default:
		return "Field " + err.Field() + " is invalid"
	}
}
