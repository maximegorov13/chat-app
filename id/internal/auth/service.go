package auth

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (string, error)
	Logout(ctx context.Context, token string) error
}
