package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/entity"
	"github.com/edwintantawi/taskit/internal/user/usecase/mocks"
)

type UserUsecaseRegisterTestSuite struct {
	suite.Suite
}

func TestUserUsecaseRegisterSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseRegisterTestSuite))
}

func (s *UserUsecaseRegisterTestSuite) TestExecute() {
	s.Run("it should return error when entity validation fail", func() {
		ctx := context.Background()
		payload := RegisterPayload{}

		usecase := NewRegister(nil, nil)
		output, err := usecase.Execute(ctx, payload)

		s.Error(err)
		s.Empty(output)
	})

	s.Run("it should return error ErrEmailAlreadyExists when user with that email is already exists", func() {
		ctx := context.Background()
		payload := RegisterPayload{
			Name:     "Gopher",
			Email:    "gopher@go.dev",
			Password: "secret_password",
		}

		mockUserRepository := &mocks.UserRepositoryFindSaver{}
		mockUserRepository.On("FindByEmail", ctx, payload.Email).Return(entity.User{}, nil)

		usecase := NewRegister(mockUserRepository, nil)
		output, err := usecase.Execute(ctx, payload)

		s.Equal(ErrEmailAlreadyExists, err)
		s.Empty(output)
	})

	s.Run("it should return error when fail to hash the password", func() {
		errOccurred := errors.New("failed hashing password")

		ctx := context.Background()
		payload := RegisterPayload{
			Name:     "Gopher",
			Email:    "gopher@go.dev",
			Password: "secret_password",
		}

		mockUserRepository := &mocks.UserRepositoryFindSaver{}
		mockUserRepository.On("FindByEmail", ctx, payload.Email).Return(entity.User{}, errors.New("user not found"))

		mockPasswordHasher := &mocks.PasswordHasher{}
		mockPasswordHasher.On("Hash", payload.Password).Return(nil, errOccurred)

		usecase := NewRegister(mockUserRepository, mockPasswordHasher)
		output, err := usecase.Execute(ctx, payload)

		s.Equal(errOccurred, err)
		s.Empty(output)
	})

	s.Run("it should return added user when success register new user", func() {
		ctx := context.Background()
		payload := RegisterPayload{
			Name:     "Gopher",
			Email:    "gopher@go.dev",
			Password: "secret_password",
		}

		mockUserRepository := &mocks.UserRepositoryFindSaver{}
		mockUserRepository.On("FindByEmail", ctx, payload.Email).Return(entity.User{}, errors.New("user not found"))
		mockUserRepository.On("Save", ctx, entity.NewUser{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: "hashed_password",
		}).Return(entity.AddedUser{ID: "user-xxxxx", Email: payload.Email}, nil)

		mockPasswordHasher := &mocks.PasswordHasher{}
		mockPasswordHasher.On("Hash", payload.Password).Return([]byte("hashed_password"), nil)

		usecase := NewRegister(mockUserRepository, mockPasswordHasher)
		output, err := usecase.Execute(ctx, payload)

		s.NoError(err)
		s.Equal("user-xxxxx", output.ID)
		s.Equal(payload.Email, output.Email)
	})
}
