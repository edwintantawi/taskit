package usecase

import (
	"context"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type usecase struct {
	userRepository domain.UserRepository
	hashProvider   domain.HashProvider
}

// New create a new user usecase.
func New(userRepository domain.UserRepository, hashProvider domain.HashProvider) domain.UserUsecase {
	return &usecase{userRepository: userRepository, hashProvider: hashProvider}
}

// Create create a new user.
func (u *usecase) Create(ctx context.Context, payload *dto.CreateUserIn) (dto.CreateUserOut, error) {
	user := &entity.User{Name: payload.Name, Email: payload.Email, Password: payload.Password}
	if err := user.Validate(); err != nil {
		return dto.CreateUserOut{}, err
	}

	if err := u.userRepository.VerifyAvailableEmail(ctx, user.Email); err != nil {
		return dto.CreateUserOut{}, err
	}

	securePassword, err := u.hashProvider.Hash(user.Password)
	if err != nil {
		return dto.CreateUserOut{}, err
	}
	user.Password = string(securePassword)

	id, err := u.userRepository.Store(ctx, user)
	if err != nil {
		return dto.CreateUserOut{}, err
	}
	return dto.CreateUserOut{ID: id, Email: user.Email}, nil
}
