package domain

import (
	"context"
	"errors"
)

var (
	ErrPasswordIncorrect = errors.New("auth.usecase.password_incorrect")
)

// UserUsecase represent user usecase contract.
type UserUsecase interface {
	Create(ctx context.Context, payload *CreateUserIn) (CreateUserOut, error)
}

// AuthUsecase represent auth usecase contract.
type AuthUsecase interface {
	Login(ctx context.Context, payload *LoginAuthIn) (LoginAuthOut, error)
	Logout(ctx context.Context, payload *LogoutAuthIn) error
	GetProfile(ctx context.Context, payload *GetProfileAuthIn) (GetProfileAuthOut, error)
	Refresh(ctx context.Context, payload *RefreshAuthIn) (RefreshAuthOut, error)
}

// TaskUsecase represent task usecase contract.
type TaskUsecase interface {
	Create(ctx context.Context, payload *CreateTaskIn) (CreateTaskOut, error)
}
