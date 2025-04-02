package http

import (
	"net/http"
	"strings"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/req"
	"github.com/maximegorov13/chat-app/id/internal/res"
)

type AuthHandlerDeps struct {
	Conf        *configs.Config
	AuthService auth.AuthService
}

type AuthHandler struct {
	conf        *configs.Config
	authService auth.AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		conf:        deps.Conf,
		authService: deps.AuthService,
	}

	router.HandleFunc("POST /api/auth/login", handler.Login())
	router.HandleFunc("POST /api/auth/logout", handler.Logout())
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[auth.LoginRequest](r)
		if err != nil {
			res.Error(w, err)
			return
		}

		token, err := h.authService.Login(r.Context(), &body.Data)
		if err != nil {
			res.Error(w, err)
			return
		}

		data := auth.LoginResponse{
			Token: token,
		}

		res.JSON(w, http.StatusOK, data, res.Meta{})
	}
}

func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			res.Error(w, apperrors.ErrUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			res.Error(w, apperrors.ErrUnauthorized)
			return
		}
		token := tokenParts[1]

		err := h.authService.Logout(r.Context(), token)
		if err != nil {
			res.Error(w, err)
			return
		}

		res.JSON(w, http.StatusOK, "Successfully logged out", res.Meta{})
	}
}
