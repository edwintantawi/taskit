package domain

import (
	"context"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// User repository errors.
var (
	ErrEmailNotAvailable = errors.New("user.repository.email_not_available")
	ErrUserNotFound      = errors.New("user.repository.user_not_found")
)

// Auth repository errors.
var (
	ErrAuthNotFound = errors.New("auth.repository.auth_not_found")
)

// UserRepository represent user repository contract.
type UserRepository interface {
	Store(ctx context.Context, u *entity.User) (entity.UserID, error)
	VerifyAvailableEmail(ctx context.Context, email string) error
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByID(ctx context.Context, id entity.UserID) (entity.User, error)
}

// AuthRepository represent auth repository contract.
type AuthRepository interface {
	Store(ctx context.Context, a *entity.Auth) error
	Delete(ctx context.Context, a *entity.Auth) error
	FindByToken(ctx context.Context, token string) (entity.Auth, error)
}

// TaskRepository represent task repository contract.
type TaskRepository interface {
	Store(ctx context.Context, t *entity.Task) (entity.TaskID, error)
	FindAllByUserID(ctx context.Context, t *entity.Task) ([]entity.Task, error)
}
