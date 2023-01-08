package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/dto"
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
	validator      *mocks.ValidatorProvider
	userRepository *mocks.UserRepository
	hashProvider   *mocks.HashProvider
}

func (s *UserUsecaseTestSuite) TestCreate() {
	type args struct {
		ctx     context.Context
		payload *dto.UserCreateIn
	}
	type expected struct {
		output dto.UserCreateOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when user validation error",
			args: args{
				ctx:     context.Background(),
				payload: &dto.UserCreateIn{},
			},
			expected: expected{
				output: dto.UserCreateOut{},
				err:    test.ErrValidator,
			},
			setup: func(d *dependency) {
				d.validator.On("Validate", mock.Anything).Return(test.ErrValidator)
			},
		},
		{
			name: "it should return error when user repository VerifyAvailableEmail return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &dto.UserCreateIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: dto.UserCreateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.validator.On("Validate", mock.Anything).Return(nil)

				d.userRepository.On("VerifyAvailableEmail", context.Background(), "gopher@go.dev").
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when hash provider Hash return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &dto.UserCreateIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: dto.UserCreateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.validator.On("Validate", mock.Anything).Return(nil)

				d.userRepository.On("VerifyAvailableEmail", context.Background(), "gopher@go.dev").
					Return(nil)

				d.hashProvider.On("Hash", "secret_password").
					Return(nil, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when user repository Store return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &dto.UserCreateIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: dto.UserCreateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.validator.On("Validate", mock.Anything).Return(nil)

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
				payload: &dto.UserCreateIn{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				output: dto.UserCreateOut{ID: "user-xxxxx", Email: "gopher@go.dev"},
				err:    nil,
			},
			setup: func(d *dependency) {
				d.validator.On("Validate", mock.Anything).Return(nil)

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
				validator:      &mocks.ValidatorProvider{},
				userRepository: &mocks.UserRepository{},
				hashProvider:   &mocks.HashProvider{},
			}
			t.setup(d)

			usecase := New(d.validator, d.userRepository, d.hashProvider)
			output, err := usecase.Create(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}
