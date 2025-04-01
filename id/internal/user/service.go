package user

import "context"

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	UpdateUser(ctx context.Context, userID int64, req *UpdateUserRequest) (*User, error)
}
