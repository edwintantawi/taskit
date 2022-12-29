package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/test"
)

type UserUsecaseTestSuite struct {
	suite.Suite
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

type dependency struct {
	userRepository *mocks.UserRepository
	hashProvider   *mocks.HashProvider
}

func (s *UserUsecaseTestSuite) TestCreate() {
	s.Run("it should return error when user validation fail", func() {
		ctx := context.Background()
		payload := &domain.CreateUserIn{}

		usecase := New(nil, nil)
		output, err := usecase.Create(ctx, payload)

		s.Error(err)
		s.Empty(output)
	})

	type args struct {
		ctx     context.Context
		payload *domain.CreateUserIn
	}
	type expected struct {
		output domain.CreateUserOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when user repository VerifyAvailableEmail return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.CreateUserOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("VerifyAvailableEmail", context.Background(), "gopher@go.dev").Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when hash provider Hash return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.CreateUserOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("VerifyAvailableEmail", context.Background(), "gopher@go.dev").Return(nil)
				d.hashProvider.On("Hash", "secret_password").Return(nil, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when user repository Store return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.CreateUserOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.userRepository.On("VerifyAvailableEmail", context.Background(), "gopher@go.dev").
					Return(nil)

				d.hashProvider.On("Hash", "secret_password").
					Return([]byte("secret_hashed_password"), nil)

				d.userRepository.On("Store", context.Background(), &entity.User{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_hashed_password"}).
					Return(entity.UserID(""), test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and output when success",
			args: args{
				ctx:     context.Background(),
				payload: &domain.CreateUserIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: domain.CreateUserOut{ID: "user-xxxxx", Email: "gopher@go.dev"},
				err:    nil,
			},
			setup: func(d *dependency) {
				d.userRepository.On("VerifyAvailableEmail", context.Background(), "gopher@go.dev").
					Return(nil)

				d.hashProvider.On("Hash", "secret_password").
					Return([]byte("secret_hashed_password"), nil)

				d.userRepository.On("Store", context.Background(), &entity.User{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_hashed_password"}).
					Return(entity.UserID("user-xxxxx"), nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				userRepository: &mocks.UserRepository{},
				hashProvider:   &mocks.HashProvider{},
			}
			t.setup(d)

			usecase := New(d.userRepository, d.hashProvider)
			output, err := usecase.Create(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}
