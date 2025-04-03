package apperrors

import (
	"net/http"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

var (
	ErrBadRequest         = NewError(http.StatusBadRequest, "bad request")
	ErrInvalidRequestBody = NewError(http.StatusBadRequest, "invalid request body")
	ErrValidationFailed   = NewError(http.StatusBadRequest, "validation failed")

	ErrUnauthorized       = NewError(http.StatusUnauthorized, "unauthorized")
	ErrInvalidCredentials = NewError(http.StatusUnauthorized, "invalid credentials")

	ErrForbidden = NewError(http.StatusForbidden, "forbidden")

	ErrNotFound = NewError(http.StatusNotFound, "not found")

	ErrUserExists = NewError(http.StatusConflict, "user already exists")

	ErrInternalServerError = NewError(http.StatusInternalServerError, "internal server error")
)
