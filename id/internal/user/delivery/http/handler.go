package http

import (
	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/user"
	"github.com/maximegorov13/chat-app/id/pkg/apperrors"
	"github.com/maximegorov13/chat-app/id/pkg/req"
	"github.com/maximegorov13/chat-app/id/pkg/res"
	"net/http"
)

type UserHandlerDeps struct {
	Conf        *configs.Config
	UserService user.UserService
}

type UserHandler struct {
	conf        *configs.Config
	userService user.UserService
}

func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		conf:        deps.Conf,
		userService: deps.UserService,
	}

	router.HandleFunc("POST /api/users", handler.Register())
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[user.RegisterRequest](r)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		u, err := h.userService.Register(r.Context(), body)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		data := user.RegisterResponse{
			ID:    u.ID,
			Login: u.Login,
			Name:  u.Name,
		}

		res.Json(w, data, http.StatusCreated)
	}
}
