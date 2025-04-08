package middleware

import (
	"net/http"
	"strconv"

	"github.com/maximegorov13/chat-app/id/internal/appcontext"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/res"
)

func CheckUserAccessByID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenUserIDStr := appcontext.GetContextUserID(r.Context())
		tokenUserID, err := strconv.ParseInt(tokenUserIDStr, 10, 64)
		if err != nil {
			res.Error(w, apperrors.ErrUnauthorized)
			return
		}

		idStr := r.PathValue("id")
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			res.Error(w, apperrors.ErrBadRequest)
			return
		}

		if tokenUserID != userID {
			res.Error(w, apperrors.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
