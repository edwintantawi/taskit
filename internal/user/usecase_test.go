package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type UsecaseTestSuite struct {
	suite.Suite
}

func TestUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (s *UsecaseTestSuite) TestCreate() {
	s.Run("it should return error when user validation fail", func() {
		ctx := context.Background()
		payload := &domain.CreateUserIn{}
		usecase := NewUsecase(nil, nil)
		r, err := usecase.Create(ctx, payload)

		s.Error(err)
		s.Empty(r)
	})

	s.Run("it should return error when verify available email error", func() {
		ctx := context.Background()
		payload := &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"}

		mockRepo := &mocks.UserRepository{}
		mockRepo.On("VerifyAvailableEmail", ctx, payload.Email).Return(ErrEmailNotAvailable)

		usecase := NewUsecase(mockRepo, nil)
		r, err := usecase.Create(ctx, payload)

		s.Equal(ErrEmailNotAvailable, err)
		s.Empty(r)
	})

	s.Run("it should return error when hashing password is error", func() {
		ctx := context.Background()
		payload := &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"}

		mockRepo := &mocks.UserRepository{}
		mockRepo.On("VerifyAvailableEmail", ctx, payload.Email).Return(nil)

		mockHash := &mocks.HashProvider{}
		mockHash.On("Hash", payload.Password).Return(nil, errors.New("password hash fail"))

		usecase := NewUsecase(mockRepo, mockHash)
		r, err := usecase.Create(ctx, payload)

		s.Equal(errors.New("password hash fail"), err)
		s.Empty(r)
	})

	s.Run("it should return error when repository create return an error", func() {
		ctx := context.Background()
		payload := &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"}

		mockRepo := &mocks.UserRepository{}
		mockRepo.On("VerifyAvailableEmail", ctx, payload.Email).Return(nil)
		mockRepo.On("Store", ctx, &entity.User{Name: payload.Name, Email: payload.Email, Password: "secure_hash_password"}).
			Return(entity.UserID(""), errors.New("repository error"))

		mockHash := &mocks.HashProvider{}
		mockHash.On("Hash", payload.Password).Return([]byte("secure_hash_password"), nil)

		usecase := NewUsecase(mockRepo, mockHash)
		r, err := usecase.Create(ctx, payload)

		s.Equal(errors.New("repository error"), err)
		s.Empty(r)
		mockRepo.AssertCalled(s.T(), "Store", ctx, &entity.User{Name: payload.Name, Email: payload.Email, Password: "secure_hash_password"})
	})

	s.Run("it should return correct result when all operations successfully", func() {
		ctx := context.Background()
		payload := &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"}

		mockRepo := &mocks.UserRepository{}
		mockRepo.On("VerifyAvailableEmail", ctx, payload.Email).Return(nil)
		mockRepo.On("Store", ctx, &entity.User{Name: payload.Name, Email: payload.Email, Password: "secure_hash_password"}).
			Return(entity.UserID("xxxxx"), nil)

		mockHash := &mocks.HashProvider{}
		mockHash.On("Hash", payload.Password).Return([]byte("secure_hash_password"), nil)

		usecase := NewUsecase(mockRepo, mockHash)
		r, err := usecase.Create(ctx, payload)

		s.NoError(err)
		s.Equal(domain.CreateUserOut{ID: "xxxxx", Email: payload.Email}, r)
	})
}
