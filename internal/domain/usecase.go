package domain

import (
	"context"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain/dto"
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
	Create(ctx context.Context, payload *dto.UserCreateIn) (dto.UserCreateOut, error)
}

// AuthUsecase represent auth usecase contract.
type AuthUsecase interface {
	Login(ctx context.Context, payload *dto.AuthLoginIn) (dto.AuthLoginOut, error)
	Logout(ctx context.Context, payload *dto.AuthLogoutIn) error
	GetProfile(ctx context.Context, payload *dto.AuthProfileIn) (dto.AuthProfileOut, error)
	Refresh(ctx context.Context, payload *dto.AuthRefreshIn) (dto.AuthRefreshOut, error)
}

// TaskUsecase represent task usecase contract.
type TaskUsecase interface {
	Create(ctx context.Context, payload *dto.TaskCreateIn) (dto.TaskCreateOut, error)
	GetAll(ctx context.Context, payload *dto.TaskGetAllIn) ([]dto.TaskGetAllOut, error)
	Remove(ctx context.Context, payload *dto.TaskRemoveIn) error
	GetByID(ctx context.Context, payload *dto.TaskGetByIDIn) (dto.TaskGetByIDOut, error)
	Update(ctx context.Context, payload *dto.TaskUpdateIn) (dto.TaskUpdateOut, error)
}
