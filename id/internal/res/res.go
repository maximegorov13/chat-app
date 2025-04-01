package res

import (
	"encoding/json"
	"net/http"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
)

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		apperrors.HandleError(w, err)
	}
}
