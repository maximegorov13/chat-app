package req

import "net/http"

type Body interface {
	Validate() error
}

func HandleBody[T Body](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil, err
	}

	if err = body.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
