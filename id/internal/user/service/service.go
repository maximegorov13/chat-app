package service

import (
	"context"
	"github.com/maximegorov13/chat-app/id/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	existedUser, err := s.repo.FindByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if existedUser != nil {
		return nil, user.ErrUserExists
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

	if err = s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
