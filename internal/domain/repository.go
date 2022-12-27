package domain

import (
	"context"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

var (
	ErrEmailNotAvailable = errors.New("user.repository.email_not_available")
	ErrUserEmailNotExist = errors.New("user.repository.user_email_not_exist")
	ErrUserIDNotExist    = errors.New("user.repository.user_id_not_exist")

	ErrAuthNotExist = errors.New("auth.repository.auth_not_exist")
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
}
