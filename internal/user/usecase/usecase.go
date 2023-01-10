package usecase

import (
	"context"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type Usecase struct {
	validator      domain.ValidatorProvider
	userRepository domain.UserRepository
	hashProvider   domain.HashProvider
}

// New create a new user usecase.
func New(validator domain.ValidatorProvider, userRepository domain.UserRepository, hashProvider domain.HashProvider) Usecase {
	return Usecase{validator: validator, userRepository: userRepository, hashProvider: hashProvider}
}

// Create create a new user.
func (u *Usecase) Create(ctx context.Context, payload *dto.UserCreateIn) (dto.UserCreateOut, error) {
	user := &entity.User{Name: payload.Name, Email: payload.Email, Password: payload.Password}
	if err := u.validator.Validate(user); err != nil {
		return dto.UserCreateOut{}, err
	}

	if err := u.userRepository.VerifyAvailableEmail(ctx, user.Email); err != nil {
		return dto.UserCreateOut{}, err
	}

	securePassword, err := u.hashProvider.Hash(user.Password)
	if err != nil {
		return dto.UserCreateOut{}, err
	}
	user.Password = string(securePassword)

	id, err := u.userRepository.Store(ctx, user)
	if err != nil {
		return dto.UserCreateOut{}, err
	}
	return dto.UserCreateOut{ID: id, Email: user.Email}, nil
}
