package http

import (
	"net/http"
	"strconv"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/appcontext"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/middleware"
	"github.com/maximegorov13/chat-app/id/internal/req"
	"github.com/maximegorov13/chat-app/id/internal/res"
	"github.com/maximegorov13/chat-app/id/internal/user"
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
	router.Handle("PUT /api/users/{id}", middleware.Auth(handler.UpdateUser(), handler.conf))
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

func (h *UserHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[user.UpdateUserRequest](r)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		idStr := r.PathValue("id")
		userId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			apperrors.HandleError(w, apperrors.ErrBadRequest)
		}

		tokenUserIDStr := appcontext.GetContextUserId(r.Context())
		tokenUserID, err := strconv.ParseInt(tokenUserIDStr, 10, 64)
		if err != nil {
			apperrors.HandleError(w, apperrors.ErrBadRequest)
		}
		if tokenUserID != userId {
			apperrors.HandleError(w, apperrors.ErrForbidden)
			return
		}

		u, err := h.userService.UpdateUser(r.Context(), tokenUserID, body)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		data := user.UpdateUserResponse{
			ID:    u.ID,
			Login: u.Login,
			Name:  u.Name,
		}

		res.Json(w, data, http.StatusOK)
	}
}
