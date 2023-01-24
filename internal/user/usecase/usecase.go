package usecase

import (
	"context"

	"github.com/edwintantawi/taskit/internal/entity"
)

// UserRepositoryFindSaver  is the interface that groups the basic find and save methods.
type UserRepositoryFindSaver interface {
	UserRepositoryFinder
	UserRepositorySaver
}

// UserRepositorySaver is the interface that wraps the basic save method.
type UserRepositorySaver interface {
	Save(ctx context.Context, newUser entity.NewUser) (entity.AddedUser, error)
}

// UserRepositoryFinder is the interface that wraps the basic find method.
type UserRepositoryFinder interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}

// PasswordHasher is an interface for hashing passwords.
type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}
