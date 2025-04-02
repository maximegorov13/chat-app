package req

import (
	"encoding/json"
	"io"
)

func Decode[T Body](body io.ReadCloser) (Request[T], error) {
	var payload Request[T]
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload, err
	}

	return payload, nil
}
