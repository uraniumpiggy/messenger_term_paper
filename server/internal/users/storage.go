package users

import "context"

type Storage interface {
	RegisterUser(context.Context, *UserRegisterRequest) error
	AuthUser(context.Context, *UserLoginRequest) (string, error)
}
