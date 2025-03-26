package user

import (
	"context"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
}
