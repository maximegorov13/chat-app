package service

import (
	"context"
	"github.com/maximegorov13/chat-app/id/internal/user"
	"github.com/maximegorov13/chat-app/id/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceDeps struct {
	UserRepo user.UserRepository
}

type UserService struct {
	userRepo user.UserRepository
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		userRepo: deps.UserRepo,
	}
}

func (s *UserService) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	existedUser, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if existedUser != nil {
		return nil, apperrors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Login:    req.Login,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	if err = s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
