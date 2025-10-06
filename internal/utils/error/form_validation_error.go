package application_error

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"justcallmesu.com/rest-api/internal/api/response"
)

func FormatValidationError(passedError error) any {
	var validationErrors validator.ValidationErrors
	var numError *strconv.NumError

	if errors.As(passedError, &validationErrors) {
		var formattedErrors = make([]response.APIValidationError, len(validationErrors))

		for index, err := range validationErrors {
			formattedErrors[index] = response.APIValidationError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
			}
		}
		return formattedErrors
	} else if errors.As(passedError, &numError) {
		return fmt.Sprintf(
			"Nilai tidak valid: '%s'. Harap periksa kembali input Anda.",
			numError.Num,
		)
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
	case "number", "numeric":
		return "Field " + err.Field() + " must be a number"
	case "oneof":
		return "Field " + err.Field() + " must be one of the following: " + err.Param()
	default:
		return "Field " + err.Field() + " is invalid"
	}
}
