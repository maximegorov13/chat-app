package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/maximegorov13/chat-app/id/pkg/jwt"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/user"
)

type AuthServiceDeps struct {
	UserRepo  user.UserRepository
	TokenRepo auth.TokenRepository
	JWT       *jwt.JWT
}

type AuthService struct {
	userRepo  user.UserRepository
	tokenRepo auth.TokenRepository
	jwt       *jwt.JWT
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		userRepo:  deps.UserRepo,
		tokenRepo: deps.TokenRepo,
		jwt:       deps.JWT,
	}
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (string, error) {
	existedUser, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return "", err
	}
	if existedUser == nil {
		return "", apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(req.Password))
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := s.jwt.GenerateToken(existedUser.ID, existedUser.Login, existedUser.Name, time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return apperrors.ErrUnauthorized
	}

	return s.tokenRepo.InvalidateToken(ctx, token, time.Hour)
}
