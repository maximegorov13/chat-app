package middleware

import (
	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/pkg/appcontext"
	"github.com/maximegorov13/chat-app/id/pkg/apperrors"
	"github.com/maximegorov13/chat-app/id/pkg/jwt"
	"net/http"
	"strings"
)

func Auth(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}
		tokenString := tokenParts[1]

		j := jwt.New(conf.Auth.Secret)
		valid, claims := j.ValidateToken(tokenString)
		if !valid {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		if j.IsTokenExpired(tokenString) {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		ctx := appcontext.SetContextUserId(r.Context(), claims.Subject)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
