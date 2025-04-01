package auth

import (
	"context"

	"github.com/maximegorov13/chat-app/id/internal/user"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*user.User, error)
}
