package user

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByLogin(ctx context.Context, login string) (*User, error)
}
