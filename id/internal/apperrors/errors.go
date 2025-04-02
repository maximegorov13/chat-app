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

func NewError(err error, code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrBadRequest         = NewError(nil, http.StatusBadRequest, "bad request")
	ErrInvalidRequestBody = NewError(nil, http.StatusBadRequest, "invalid request body")
	ErrValidationFailed   = NewError(nil, http.StatusBadRequest, "validation failed")

	ErrUnauthorized       = NewError(nil, http.StatusUnauthorized, "unauthorized")
	ErrInvalidCredentials = NewError(nil, http.StatusUnauthorized, "invalid credentials")

	ErrForbidden = NewError(nil, http.StatusForbidden, "forbidden")

	ErrNotFound = NewError(nil, http.StatusNotFound, "not found")

	ErrUserExists = NewError(nil, http.StatusConflict, "user already exists")

	ErrInternalServerError = NewError(nil, http.StatusInternalServerError, "internal server error")
)
