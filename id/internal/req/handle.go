package req

import (
	"net/http"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
)

func HandleBody[T Body](r *http.Request) (*Request[T], error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		return nil, apperrors.ErrInvalidRequestBody
	}

	if err = body.Data.Validate(); err != nil {
		return nil, apperrors.NewError(err, apperrors.ErrValidationFailed.Code, apperrors.ErrValidationFailed.Message)
	}

	return &body, nil
}
