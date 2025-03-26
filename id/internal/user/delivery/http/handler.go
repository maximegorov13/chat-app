package http

import (
	"errors"
	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/user"
	"github.com/maximegorov13/chat-app/id/pkg/req"
	"github.com/maximegorov13/chat-app/id/pkg/res"
	"net/http"
)

type UserHandlerDeps struct {
	Conf        *configs.Config
	UserService user.UserService
}

type UserHandler struct {
	Conf        *configs.Config
	UserService user.UserService
}

func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		Conf:        deps.Conf,
		UserService: deps.UserService,
	}

	router.HandleFunc("POST /api/users", handler.Register())
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[user.RegisterRequest](w, r)
		if err != nil {
			return
		}

		u, err := h.UserService.Register(r.Context(), body)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrUserExists):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := user.RegisterResponse{
			ID:    u.ID,
			Login: u.Login,
			Name:  u.Name,
		}

		res.Json(w, response, http.StatusCreated)
	}
}
