package req

import (
	"net/http"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
)

type Body interface {
	Validate() error
}

func HandleBody[T Body](r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		return nil, apperrors.ErrInvalidRequestBody
	}

	if err = body.Validate(); err != nil {
		return nil, apperrors.New(err, apperrors.ErrValidationFailed.Code, apperrors.ErrValidationFailed.Message)
	}

	return &body, nil
}
