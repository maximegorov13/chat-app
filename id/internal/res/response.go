package res

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-ozzo/ozzo-validation"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
)

func JSON(w http.ResponseWriter, statusCode int, data any, meta ResponseMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Data: data,
		Meta: meta,
	})
}

func Error(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		sendAppError(w, appErr)
		return
	}

	var valErr validation.Errors
	if errors.As(err, &valErr) {
		sendValidateError(w, valErr)
		return
	}

	sendDefaultError(w)
}

func sendAppError(w http.ResponseWriter, err *apperrors.Error) {
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(Response{
		Error: &ErrorResponse{
			Code:    err.Code,
			Message: err.Message,
		},
	})
}

func sendValidateError(w http.ResponseWriter, errs validation.Errors) {
	w.WriteHeader(apperrors.ErrValidationFailed.Code)

	details := make([]ErrorDetail, 0, len(errs))
	for field, err := range errs {
		details = append(details, ErrorDetail{
			Field:   field,
			Message: err.Error(),
		})
	}

	json.NewEncoder(w).Encode(Response{
		Error: &ErrorResponse{
			Code:    apperrors.ErrValidationFailed.Code,
			Message: apperrors.ErrValidationFailed.Message,
			Details: details,
		},
	})
}

func sendDefaultError(w http.ResponseWriter) {
	w.WriteHeader(apperrors.ErrInternalServerError.Code)
	json.NewEncoder(w).Encode(Response{
		Error: &ErrorResponse{
			Code:    apperrors.ErrInternalServerError.Code,
			Message: apperrors.ErrInternalServerError.Message,
		},
	})
}
