package domain

import (
	"context"

	"github.com/edwintantawi/taskit/internal/user/domain/entity"
)

// UserRepository represent user repository contract
type UserRepository interface {
	Store(ctx context.Context, u *entity.User) (entity.UserID, error)
	VerifyAvailableEmail(ctx context.Context, email string) error
}
