package req

import (
	"github.com/maximegorov13/chat-app/id/pkg/apperrors"
	"net/http"
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
		return nil, apperrors.New(err, http.StatusBadRequest, "invalid request body")
	}

	return &body, nil
}
