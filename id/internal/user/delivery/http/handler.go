package http

import (
	"net/http"
	"strconv"

	"github.com/maximegorov13/chat-app/id/pkg/jwt"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/middleware"
	"github.com/maximegorov13/chat-app/id/internal/req"
	"github.com/maximegorov13/chat-app/id/internal/res"
	"github.com/maximegorov13/chat-app/id/internal/user"
)

type UserHandlerDeps struct {
	Conf        *configs.Config
	UserService user.UserService
	TokenRepo   auth.TokenRepository
	JWT         *jwt.JWT
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
	router.Handle("PUT /api/users/{id}", middleware.Auth(middleware.CheckUserAccessByID(handler.UpdateUser()), middleware.AuthDeps{
		Conf:      deps.Conf,
		TokenRepo: deps.TokenRepo,
		JWT:       deps.JWT,
	}))
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[user.RegisterRequest](r)
		if err != nil {
			res.Error(w, err)
			return
		}

		u, err := h.userService.Register(r.Context(), &body.Data)
		if err != nil {
			res.Error(w, err)
			return
		}

		data := user.RegisterResponse{
			ID:    u.ID,
			Login: u.Login,
			Name:  u.Name,
		}

		res.JSON(w, http.StatusCreated, data, nil)
	}
}

func (h *UserHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[user.UpdateUserRequest](r)
		if err != nil {
			res.Error(w, err)
			return
		}

		idStr := r.PathValue("id")
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			res.Error(w, apperrors.ErrBadRequest)
			return
		}

		u, err := h.userService.UpdateUser(r.Context(), userID, &body.Data)
		if err != nil {
			res.Error(w, err)
			return
		}

		data := user.UpdateUserResponse{
			ID:    u.ID,
			Login: u.Login,
			Name:  u.Name,
		}

		res.JSON(w, http.StatusOK, data, nil)
	}
}
