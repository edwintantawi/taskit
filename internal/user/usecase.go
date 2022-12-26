package user

import (
	"context"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type usecase struct {
	userRepository domain.UserRepository
	hashProvider   domain.HashProvider
}

func NewUsecase(userRepository domain.UserRepository, hashProvider domain.HashProvider) domain.UserUsecase {
	return &usecase{userRepository: userRepository, hashProvider: hashProvider}
}

// Create create a new user.
func (u *usecase) Create(ctx context.Context, payload *domain.CreateUserIn) (domain.CreateUserOut, error) {
	user := &entity.User{Name: payload.Name, Email: payload.Email, Password: payload.Password}
	if err := user.Validate(); err != nil {
		return domain.CreateUserOut{}, err
	}

	if err := u.userRepository.VerifyAvailableEmail(ctx, user.Email); err != nil {
		return domain.CreateUserOut{}, err
	}

	securePassword, err := u.hashProvider.Hash(user.Password)
	if err != nil {
		return domain.CreateUserOut{}, err
	}
	user.Password = string(securePassword)

	id, err := u.userRepository.Store(ctx, user)
	if err != nil {
		return domain.CreateUserOut{}, err
	}
	return domain.CreateUserOut{ID: id, Email: user.Email}, nil
}
