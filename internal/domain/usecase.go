package domain

import (
	"context"
	"errors"
)

// Auth usecase errors.
var (
	ErrEmailNotExist     = errors.New("auth.usecase.email_not_exist")
	ErrPasswordIncorrect = errors.New("auth.usecase.password_incorrect")
)

// Task usecase errors.
var (
	ErrTaskAuthorization = errors.New("task.usecase.task_forbidden")
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
	GetAll(ctx context.Context, payload *GetAllTaskIn) ([]GetAllTaskOut, error)
	Remove(ctx context.Context, payload *RemoveTaskIn) error
	GetByID(ctx context.Context, payload *GetTaskByIDIn) (GetTaskByIDOut, error)
	Update(ctx context.Context, payload *UpdateTaskIn) (UpdateTaskOut, error)
}
