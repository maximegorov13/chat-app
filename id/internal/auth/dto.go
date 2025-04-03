package auth

import "github.com/go-ozzo/ozzo-validation"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Login, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type LoginResponse struct {
	Token string `json:"token"`
}
