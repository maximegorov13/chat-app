package res

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
)

func JSON(w http.ResponseWriter, statusCode int, data any, meta Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Data: data,
		Meta: meta,
	})
}

func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var e *apperrors.Error
	if errors.As(err, &e) {
		w.WriteHeader(e.Code)
		json.NewEncoder(w).Encode(Response{
			Error: &ErrorDetails{
				Code:    e.Code,
				Message: e.Message,
			},
		})
	} else {
		w.WriteHeader(apperrors.ErrInternalServerError.Code)
		json.NewEncoder(w).Encode(Response{
			Error: &ErrorDetails{
				Code:    apperrors.ErrInternalServerError.Code,
				Message: apperrors.ErrInternalServerError.Message,
			},
		})
	}
}
