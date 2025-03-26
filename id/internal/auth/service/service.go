package service

import (
	"context"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/user"
)

type AuthServiceDeps struct {
	UserRepo user.UserRepository
}

type AuthService struct {
	userRepo user.UserRepository
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		userRepo: deps.UserRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*user.User, error) {
}
