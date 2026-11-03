package request

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"justcallmesu.com/golang-base/internal/api/response"
)

func formatValidationError(errs validator.ValidationErrors) []response.APIValidationError {
	output := make([]response.APIValidationError, len(errs))

	for index, fieldError := range errs {
		output[index] = response.APIValidationError{
			Field:   fieldError.Field(),
			Message: messageFor(fieldError),
		}
	}
	return output
}

func messageFor(err validator.FieldError) string {

	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		if err.Kind() == reflect.String {
			return err.Field() + " must be at least " + err.Param() + " characters long"
		}
		return err.Field() + " must be at least " + err.Param()
	case "max":
		if err.Kind() == reflect.String {
			return err.Field() + " must be at most " + err.Param() + " characters long"
		}
		return err.Field() + " must be at most " + err.Param()
	case "len":
		return err.Field() + " must be exactly " + err.Param() + " characters long"
	case "oneof":
		return err.Field() + " must be one of the following values: " + err.Param()
	case "uuid":
		return err.Field() + " must be a valid UUID"
	case "gt":
		return err.Field() + " must be greater than " + err.Param()
	case "gte":
		return err.Field() + " must be greater than or equal to " + err.Param()
	case "lt":
		return err.Field() + " must be less than " + err.Param()
	case "lte":
		return err.Field() + " must be less than or equal to " + err.Param()
	case "numeric":
		return err.Field() + " must be a numeric value"
	case "alphanum":
		return err.Field() + " must contain only alphanumeric characters"
	case "alpha":
		return err.Field() + " must contain only alphabetic characters"
	case "url":
		return err.Field() + " must be a valid URL"
	case "base64":
		return err.Field() + " must be a valid base64 encoded string"
	case "datetime":
		return err.Field() + " must be a valid datetime string"
	case "json":
		return err.Field() + " must be a valid JSON string"
	case "eqfield":
		return err.Field() + " must be equal to " + err.Param()
	case "nefield":
		return err.Field() + " must not be equal to " + err.Param()
	case "isdefault":
		return err.Field() + " must be the default value"
	case "isnotzero":
		return err.Field() + " must not be zero"
	case "isnotpositive":
		return err.Field() + " must not be positive"
	case "isnotinfinite":
		return err.Field() + " must not be infinite"
	case "isnotnan":
		return err.Field() + " must not be NaN"
	case "isnotemptystring":
		return err.Field() + " must not be an empty string"
	case "isnotwhitespace":
		return err.Field() + " must not be whitespace"
	case "isnotblankstring":
		return err.Field() + " must not be a blank string"
	default:
		return err.Field() + " is invalid"
	}
}
