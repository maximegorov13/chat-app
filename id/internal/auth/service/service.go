package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/user"
)

type AuthServiceDeps struct {
	UserRepo  user.UserRepository
	TokenRepo auth.TokenRepository
}

type AuthService struct {
	userRepo  user.UserRepository
	tokenRepo auth.TokenRepository
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		userRepo:  deps.UserRepo,
		tokenRepo: deps.TokenRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*user.User, error) {
	existedUser, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if existedUser == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(req.Password))
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	return existedUser, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return apperrors.ErrUnauthorized
	}

	return s.tokenRepo.InvalidateToken(ctx, token, time.Hour)
}
