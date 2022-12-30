package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/test"
)

type AuthUsecaseTestSuite struct {
	suite.Suite
}

func TestAuthUsecaseSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}

type dependency struct {
	authRepository *mocks.AuthRepository
	userRepository *mocks.UserRepository
	hashProvider   *mocks.HashProvider
	jwtProvider    *mocks.JWTProvider
}

func (s *AuthUsecaseTestSuite) TestLogin() {
	type args struct {
		ctx     context.Context
		payload *domain.LoginAuthIn
	}
	type expected struct {
		output domain.LoginAuthOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error ErrEmailNotExist when email not found",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev"},
			},
			expected: expected{
				output: domain.LoginAuthOut{},
				err:    domain.ErrEmailNotExist,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{}, domain.ErrUserNotFound)
			},
		},
		{
			name: "it should return error when user repository FindByEmail return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev"},
			},
			expected: expected{
				output: domain.LoginAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error ErrPasswordIncorrect when password is incorrect",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.LoginAuthOut{},
				err:    domain.ErrPasswordIncorrect,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{Password: "secret_hashed_password"}, nil)

				d.hashProvider.On("Compare", "secret_password", "secret_hashed_password").
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when generate access token failed",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.LoginAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{ID: "user-xxxxx", Password: "secret_hashed_password"}, nil)

				d.hashProvider.On("Compare", "secret_password", "secret_hashed_password").
					Return(nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("", time.Time{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when generate refresh token failed",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.LoginAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{ID: "user-xxxxx", Password: "secret_hashed_password"}, nil)

				d.hashProvider.On("Compare", "secret_password", "secret_hashed_password").
					Return(nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("", time.Time{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when auth respository Store return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.LoginAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{ID: "user-xxxxx", Password: "secret_hashed_password"}, nil)

				d.hashProvider.On("Compare", "secret_password", "secret_hashed_password").
					Return(nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("yyyyy.yyyyy.yyyyy", test.TimeAfterNow, nil)

				d.authRepository.On("Store", context.Background(), &entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and output when success",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LoginAuthIn{Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.LoginAuthOut{
					AccessToken:  "xxxxx.xxxxx.xxxxx",
					RefreshToken: "yyyyy.yyyyy.yyyyy",
				},
				err: nil,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByEmail", context.Background(), "gopher@go.dev").
					Return(entity.User{ID: "user-xxxxx", Password: "secret_hashed_password"}, nil)

				d.hashProvider.On("Compare", "secret_password", "secret_hashed_password").
					Return(nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("yyyyy.yyyyy.yyyyy", test.TimeAfterNow, nil)

				d.authRepository.On("Store", context.Background(), &entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				userRepository: &mocks.UserRepository{},
				authRepository: &mocks.AuthRepository{},
				hashProvider:   &mocks.HashProvider{},
				jwtProvider:    &mocks.JWTProvider{},
			}
			t.setup(d)

			usecase := New(d.authRepository, d.userRepository, d.hashProvider, d.jwtProvider)
			output, err := usecase.Login(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}

func (s *AuthUsecaseTestSuite) TestLogout() {
	s.Run("it should return error when validation fail", func() {
		ctx := context.Background()
		payload := &domain.LogoutAuthIn{}
		usecase := New(nil, nil, nil, nil)

		err := usecase.Logout(ctx, payload)

		s.Error(err)
	})

	type args struct {
		ctx     context.Context
		payload *domain.LogoutAuthIn
	}
	type expected struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when auth VerifyAvailableByToken return error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LogoutAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				err: test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("VerifyAvailableByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when auth Delete repository return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LogoutAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				err: test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("VerifyAvailableByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(nil)

				d.authRepository.On("Delete", context.Background(), &entity.Auth{Token: "yyyyy.yyyyy.yyyyy"}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil when successfully delete authentication",
			args: args{
				ctx:     context.Background(),
				payload: &domain.LogoutAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				err: nil,
			},
			setup: func(d *dependency) {
				d.authRepository.On("VerifyAvailableByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(nil)

				d.authRepository.On("Delete", context.Background(), &entity.Auth{Token: "yyyyy.yyyyy.yyyyy"}).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				authRepository: &mocks.AuthRepository{},
			}
			t.setup(d)

			usecase := New(d.authRepository, d.userRepository, d.hashProvider, d.jwtProvider)
			err := usecase.Logout(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
		})
	}
}

func (s *AuthUsecaseTestSuite) TestGetProfile() {
	type args struct {
		ctx     context.Context
		payload *domain.GetProfileAuthIn
	}
	type expected struct {
		output domain.GetProfileAuthOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when user repository FindByID return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.GetProfileAuthIn{UserID: "user-xxxxx"},
			},
			expected: expected{
				output: domain.GetProfileAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByID", context.Background(), entity.UserID("user-xxxxx")).
					Return(entity.User{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and output when success",
			args: args{
				ctx:     context.Background(),
				payload: &domain.GetProfileAuthIn{UserID: "user-xxxxx"},
			},
			expected: expected{
				output: domain.GetProfileAuthOut{
					ID:    "user-xxxxx",
					Name:  "Gopher",
					Email: "gopher@go.dev",
				},
				err: nil,
			},
			setup: func(d *dependency) {
				d.userRepository.On("FindByID", context.Background(), entity.UserID("user-xxxxx")).
					Return(entity.User{ID: entity.UserID("user-xxxxx"), Name: "Gopher", Email: "gopher@go.dev"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				userRepository: &mocks.UserRepository{},
			}
			t.setup(d)

			usecase := New(d.authRepository, d.userRepository, d.hashProvider, d.jwtProvider)
			output, err := usecase.GetProfile(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}

func (s *AuthUsecaseTestSuite) TestRefresh() {
	type args struct {
		ctx     context.Context
		payload *domain.RefreshAuthIn
	}
	type expected struct {
		output domain.RefreshAuthOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when auth repository FindByToken return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error ErrAuthTokenExpired when token is expired",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{},
				err:    entity.ErrAuthTokenExpired,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{ExpiresAt: test.TimeBeforeNow}, nil)
			},
		},
		{
			name: "it should return error when generate new access token failed",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}, nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("", time.Time{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when generate new refresh token failed",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}, nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("", time.Time{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when auth respository Delete return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}, nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("zzzzz.zzzzz.zzzzz", test.TimeAfterNow, nil)

				d.authRepository.On("Delete", context.Background(), &entity.Auth{Token: "yyyyy.yyyyy.yyyyy", UserID: "user-xxxxx", ExpiresAt: test.TimeAfterNow}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when auth respository Store return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}, nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("zzzzz.zzzzz.zzzzz", test.TimeAfterNow, nil)

				d.authRepository.On("Delete", context.Background(), &entity.Auth{Token: "yyyyy.yyyyy.yyyyy", UserID: "user-xxxxx", ExpiresAt: test.TimeAfterNow}).
					Return(nil)

				d.authRepository.On("Store", context.Background(), &entity.Auth{UserID: "user-xxxxx", Token: "zzzzz.zzzzz.zzzzz", ExpiresAt: test.TimeAfterNow}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and output when success",
			args: args{
				ctx:     context.Background(),
				payload: &domain.RefreshAuthIn{RefreshToken: "yyyyy.yyyyy.yyyyy"},
			},
			expected: expected{
				output: domain.RefreshAuthOut{
					AccessToken:  "xxxxx.xxxxx.xxxxx",
					RefreshToken: "zzzzz.zzzzz.zzzzz",
				},
				err: nil,
			},
			setup: func(d *dependency) {
				d.authRepository.On("FindByToken", context.Background(), "yyyyy.yyyyy.yyyyy").
					Return(entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}, nil)

				d.jwtProvider.On("GenerateAccessToken", entity.UserID("user-xxxxx")).
					Return("xxxxx.xxxxx.xxxxx", test.TimeAfterNow, nil)

				d.jwtProvider.On("GenerateRefreshToken", entity.UserID("user-xxxxx")).
					Return("zzzzz.zzzzz.zzzzz", test.TimeAfterNow, nil)

				d.authRepository.On("Delete", context.Background(), &entity.Auth{Token: "yyyyy.yyyyy.yyyyy", UserID: "user-xxxxx", ExpiresAt: test.TimeAfterNow}).
					Return(nil)

				d.authRepository.On("Store", context.Background(), &entity.Auth{UserID: "user-xxxxx", Token: "zzzzz.zzzzz.zzzzz", ExpiresAt: test.TimeAfterNow}).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				authRepository: &mocks.AuthRepository{},
				jwtProvider:    &mocks.JWTProvider{},
			}
			t.setup(d)

			usecase := New(d.authRepository, d.userRepository, d.hashProvider, d.jwtProvider)
			output, err := usecase.Refresh(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}
