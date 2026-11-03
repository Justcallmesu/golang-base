package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"justcallmesu.com/golang-base/internal/api/response"
	"justcallmesu.com/golang-base/internal/common/apperror"
)

func BindError(err error) error {
	var validationErrors validator.ValidationErrors
	var typeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError

	switch {
	case errors.As(err, &validationErrors):
		return apperror.Validation(formatValidationError(validationErrors))
	case errors.As(err, &typeError):
		return apperror.Validation([]response.APIValidationError{
			{
				Field:   typeError.Field,
				Message: fmt.Sprintf("%s must be a %s", typeError.Field, typeError.Type.Kind()),
			},
		})
	case errors.As(err, &syntaxError), errors.Is(err, io.ErrUnexpectedEOF):
		return apperror.BadRequest("Invalid JSON body")
	case errors.Is(err, io.EOF):
		return apperror.BadRequest("Empty request body")
	default:
		return apperror.BadRequest("Invalid request body")
	}
}
