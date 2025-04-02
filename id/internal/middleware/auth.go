package middleware

import (
	"net/http"
	"strings"

	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/pkg/jwt"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/appcontext"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
)

type AuthDeps struct {
	Conf      *configs.Config
	TokenRepo auth.TokenRepository
}

func Auth(next http.Handler, deps AuthDeps) http.Handler {
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
		token := tokenParts[1]

		invalid, err := deps.TokenRepo.IsTokenInvalid(r.Context(), token)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
		if invalid {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		j := jwt.New(deps.Conf.Auth.Secret)
		valid, claims := j.ValidateToken(token)
		if !valid {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		if j.IsTokenExpired(token) {
			apperrors.HandleError(w, apperrors.ErrUnauthorized)
			return
		}

		ctx := appcontext.SetContextUserId(r.Context(), claims.Subject)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
