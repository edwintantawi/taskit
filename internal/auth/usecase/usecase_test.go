package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type AuthUsecaseTestSuite struct {
	suite.Suite
}

func TestAuthUsecaseSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}

func (s *AuthUsecaseTestSuite) TestLogin() {
	s.Run("it should return an error if the user does not exist", func() {
		ctx := context.Background()
		payload := &domain.LoginAuthIn{Email: "gopher@go.dev"}

		mockUserRepo := &mocks.UserRepository{}
		mockUserRepo.On("FindByEmail", ctx, payload.Email).Return(entity.User{}, domain.ErrUserEmailNotFound)

		usecase := New(nil, mockUserRepo, nil, nil)
		_, err := usecase.Login(ctx, payload)

		s.Equal(domain.ErrUserEmailNotFound, err)
	})

	s.Run("it should return an error if the password is incorrect", func() {
		ctx := context.Background()
		payload := &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"}

		mockUserRepo := &mocks.UserRepository{}
		mockUserRepo.On("FindByEmail", ctx, payload.Email).Return(entity.User{Password: "secret_password"}, nil)

		mockHashProvider := &mocks.HashProvider{}
		mockHashProvider.On("Compare", payload.Password, mock.Anything).Return(domain.ErrPasswordIncorrect)

		usecase := New(nil, mockUserRepo, mockHashProvider, nil)
		_, err := usecase.Login(ctx, payload)

		s.Equal(domain.ErrPasswordIncorrect, err)
	})

	s.Run("it should return an error if the access token cannot be generated", func() {
		ctx := context.Background()
		payload := &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"}
		auth := &entity.Auth{UserID: "xxxxx"}

		mockUserRepo := &mocks.UserRepository{}
		mockUserRepo.On("FindByEmail", ctx, payload.Email).Return(entity.User{ID: "xxxxx"}, nil)

		mockHashProvider := &mocks.HashProvider{}
		mockHashProvider.On("Compare", payload.Password, mock.Anything).Return(nil)

		mockJWTProvider := &mocks.JWTProvider{}
		mockJWTProvider.On("GenerateAccessToken", auth.UserID).Return("", time.Time{}, errors.New("failed to generate token"))

		usecase := New(nil, mockUserRepo, mockHashProvider, mockJWTProvider)
		_, err := usecase.Login(ctx, payload)

		s.Equal(errors.New("failed to generate token"), err)
	})

	s.Run("it should return an error if the refresh token cannot be generated", func() {
		ctx := context.Background()
		payload := &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"}
		auth := &entity.Auth{UserID: "xxxxx"}

		mockUserRepo := &mocks.UserRepository{}
		mockUserRepo.On("FindByEmail", ctx, payload.Email).Return(entity.User{ID: "xxxxx"}, nil)

		mockHashProvider := &mocks.HashProvider{}
		mockHashProvider.On("Compare", payload.Password, mock.Anything).Return(nil)

		mockJWTProvider := &mocks.JWTProvider{}
		mockJWTProvider.On("GenerateAccessToken", auth.UserID).Return("", time.Time{}, nil)
		mockJWTProvider.On("GenerateRefreshToken", auth.UserID).Return("", time.Time{}, errors.New("failed to generate token"))

		usecase := New(nil, mockUserRepo, mockHashProvider, mockJWTProvider)
		_, err := usecase.Login(ctx, payload)

		s.Equal(errors.New("failed to generate token"), err)
	})

	s.Run("it should return an error if the auth cannot be saved", func() {
		ctx := context.Background()
		payload := &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"}
		auth := &entity.Auth{UserID: "xxxxx"}

		mockUserRepo := &mocks.UserRepository{}
		mockUserRepo.On("FindByEmail", ctx, payload.Email).Return(entity.User{ID: "xxxxx"}, nil)

		mockHashProvider := &mocks.HashProvider{}
		mockHashProvider.On("Compare", payload.Password, mock.Anything).Return(nil)

		mockJWTProvider := &mocks.JWTProvider{}
		mockJWTProvider.On("GenerateAccessToken", auth.UserID).Return("", time.Time{}, nil)
		mockJWTProvider.On("GenerateRefreshToken", auth.UserID).Return("", time.Time{}, nil)

		mockAuthRepo := &mocks.AuthRepository{}
		mockAuthRepo.On("Store", mock.Anything, mock.Anything).Return(errors.New("failed to save auth"))

		usecase := New(mockAuthRepo, mockUserRepo, mockHashProvider, mockJWTProvider)
		_, err := usecase.Login(ctx, payload)

		s.Equal(errors.New("failed to save auth"), err)
	})

	s.Run("it should return an auth if the login is successful", func() {
		ctx := context.Background()
		payload := &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"}
		auth := &entity.Auth{UserID: "xxxxx", Token: "refresh_token", ExpiresAt: time.Now().Add(time.Hour * 24 * 7)}

		mockUserRepo := &mocks.UserRepository{}
		mockUserRepo.On("FindByEmail", ctx, payload.Email).Return(entity.User{ID: "xxxxx"}, nil)

		mockHashProvider := &mocks.HashProvider{}
		mockHashProvider.On("Compare", payload.Password, mock.Anything).Return(nil)

		mockJWTProvider := &mocks.JWTProvider{}
		mockJWTProvider.On("GenerateAccessToken", auth.UserID).Return("access_token", time.Time{}, nil)
		mockJWTProvider.On("GenerateRefreshToken", auth.UserID).Return("refresh_token", auth.ExpiresAt, nil)

		mockAuthRepo := &mocks.AuthRepository{}
		mockAuthRepo.On("Store", ctx, auth).Return(nil)

		usecase := New(mockAuthRepo, mockUserRepo, mockHashProvider, mockJWTProvider)
		r, err := usecase.Login(ctx, payload)

		s.NoError(err)
		s.Equal("access_token", r.AccessToken)
		s.Equal("refresh_token", r.RefreshToken)
	})
}