package user

import "github.com/go-ozzo/ozzo-validation"

type RegisterRequest struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Login, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 30)),
	)
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type UpdateUserRequest struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Login, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 30)),
	)
}

type UpdateUserResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}
