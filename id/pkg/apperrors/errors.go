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
	ErrInvalidRequestBody = New(nil, http.StatusBadRequest, "invalid request body")
	ErrUserExists         = New(nil, http.StatusConflict, "user already exists")

	ErrInternalServerError = New(nil, http.StatusInternalServerError, "internal server error")
)
