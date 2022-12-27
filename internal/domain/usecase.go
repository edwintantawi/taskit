package domain

import (
	"context"
	"errors"
)

var (
	ErrPasswordIncorrect = errors.New("auth.usecase.password_incorrect")
)

type UserUsecase interface {
	Create(ctx context.Context, payload *CreateUserIn) (CreateUserOut, error)
}

// AuthUsecase represent auth usecase contract.
type AuthUsecase interface {
	Login(ctx context.Context, payload *LoginAuthIn) (LoginAuthOut, error)
	Logout(ctx context.Context, payload *LogoutAuthIn) error
}
