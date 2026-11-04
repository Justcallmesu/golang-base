package apperror

import "errors"

type Kind int

const (
	//50x
	KindInternal Kind = iota

	// 40x
	KindBadRequest
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
	KindUnprocessableEntity
)

type Error struct {
	Kind    Kind
	Message string
	Details any
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

func NotFound(message string) *Error {
	return &Error{
		Kind:    KindNotFound,
		Message: message,
	}
}

func BadRequest(message string) *Error {
	return &Error{
		Kind:    KindBadRequest,
		Message: message,
	}
}

func Unauthorized(message string) *Error {
	return &Error{
		Kind:    KindUnauthorized,
		Message: message,
	}
}

func Forbidden(message string) *Error {
	return &Error{
		Kind:    KindForbidden,
		Message: message,
	}
}

func Conflict(message string) *Error {
	return &Error{
		Kind:    KindConflict,
		Message: message,
	}
}

func UnprocessableEntity(message string) *Error {
	return &Error{
		Kind:    KindUnprocessableEntity,
		Message: message,
	}
}

func Internal(cause error) *Error {
	return &Error{
		Kind:    KindInternal,
		Message: "Internal Server Error",
		Err:     cause,
	}
}

func Validation(details any) *Error {
	return &Error{
		Kind:    KindBadRequest,
		Message: "Validation Error",
		Details: details,
	}
}

func IsKind(err error, kind Kind) bool {
	var appErr *Error
	if ok := errors.As(err, &appErr); ok {
		return appErr.Kind == kind
	}
	return false
}
