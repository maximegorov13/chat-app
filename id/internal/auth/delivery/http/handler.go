package http

import (
	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/pkg/apperrors"
	"github.com/maximegorov13/chat-app/id/pkg/jwt"
	"github.com/maximegorov13/chat-app/id/pkg/req"
	"github.com/maximegorov13/chat-app/id/pkg/res"
	"net/http"
	"time"
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
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[auth.LoginRequest](r)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		u, err := h.authService.Login(r.Context(), body)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		t, err := jwt.New(h.conf.Auth.Secret).GenerateToken(u.ID, u.Login, u.Name, time.Hour)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		data := auth.LoginResponse{
			Token: t,
		}

		res.Json(w, data, http.StatusOK)
	}
}
