package apperrors

import (
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	var e *Error
	if errors.As(err, &e) {
		http.Error(w, e.Error(), e.Code)
	} else {
		http.Error(w, ErrInternalServerError.Error(), ErrInternalServerError.Code)
	}
}
