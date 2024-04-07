package errs

import (
	"net/http"
)

// for business level

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotfoundError(message string) error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnexpectedError() error {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}

func NewValidationError(message string) error {
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}
