package usecase

import (
	"context"
	"errors"

	"github.com/edwintantawi/taskit/internal/entity"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type RegisterPayload struct {
	Name     string
	Email    string
	Password string
}

type Register struct {
	userRepository UserRepositoryFindSaver
	passwordHasher PasswordHasher
}

func NewRegister(userRepository UserRepositoryFindSaver, passwordHasher PasswordHasher) Register {
	return Register{userRepository: userRepository, passwordHasher: passwordHasher}
}

// Execute execute the register usecase.
func (r Register) Execute(ctx context.Context, payload RegisterPayload) (entity.AddedUser, error) {
	newUser := entity.NewUser(payload)
	if err := newUser.Validate(); err != nil {
		return entity.AddedUser{}, err
	}

	_, err := r.userRepository.FindByEmail(ctx, newUser.Email)
	if err == nil {
		return entity.AddedUser{}, ErrEmailAlreadyExists
	}

	hashedPassword, err := r.passwordHasher.Hash(newUser.Password)
	if err != nil {
		return entity.AddedUser{}, err
	}
	newUser.Password = string(hashedPassword)

	return r.userRepository.Save(ctx, newUser)
}
