package user

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByLogin(ctx context.Context, login string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
}
