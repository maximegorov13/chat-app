package apperrors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func New(err error, code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrBadRequest         = New(nil, http.StatusBadRequest, "bad request")
	ErrInvalidRequestBody = New(nil, http.StatusBadRequest, "invalid request body")
	ErrValidationFailed   = New(nil, http.StatusBadRequest, "validation failed")

	ErrUnauthorized       = New(nil, http.StatusUnauthorized, "unauthorized")
	ErrInvalidCredentials = New(nil, http.StatusUnauthorized, "invalid credentials")

	ErrForbidden = New(nil, http.StatusForbidden, "forbidden")

	ErrNotFound = New(nil, http.StatusNotFound, "not found")

	ErrUserExists = New(nil, http.StatusConflict, "user already exists")

	ErrInternalServerError = New(nil, http.StatusInternalServerError, "internal server error")
)
