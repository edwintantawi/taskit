package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/test"
)

type UserRepositoryTestSuite struct {
	suite.Suite
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

type dependency struct {
	mockDB     sqlmock.Sqlmock
	idProvider *mocks.IDProvider
}

func (s *UserRepositoryTestSuite) TestCreate() {
	type args struct {
		ctx  context.Context
		user *entity.User
	}
	type expected struct {
		userID entity.UserID
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when database fail to store",
			args: args{
				ctx:  context.Background(),
				user: &entity.User{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				userID: "",
				err:    test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return("user-xxxxx")
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`)).
					WithArgs("user-xxxxx", "Gopher", "gopher@go.dev", "secret_password").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error nil and user id when successfully store",
			args: args{
				ctx:  context.Background(),
				user: &entity.User{Name: "Gopher", Email: "gopher@go.dev", Password: "secret_password"},
			},
			expected: expected{
				userID: "user-xxxxx",
				err:    nil,
			},
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return("user-xxxxx")
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`)).
					WithArgs("user-xxxxx", "Gopher", "gopher@go.dev", "secret_password").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
			}

			deps := &dependency{
				mockDB:     mockDB,
				idProvider: &mocks.IDProvider{},
			}
			t.setup(deps)

			repository := New(db, deps.idProvider)
			userID, err := repository.Store(t.args.ctx, t.args.user)

			s.Equal(t.expected.userID, userID)
			s.Equal(t.expected.err, err)
		})
	}
}

func (s *UserRepositoryTestSuite) TestVerifyAvailableEmail() {
	type args struct {
		ctx   context.Context
		email string
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
			name: "it should return error when database fail to query",
			args: args{
				ctx:   context.Background(),
				email: "gopher@go.dev",
			},
			expected: expected{
				err: test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE email = $1")).
					WithArgs("gopher@go.dev").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error ErrEmailNotAvailable when email is already exist in database",
			args: args{
				ctx:   context.Background(),
				email: "gopher@go.dev",
			},
			expected: expected{
				err: domain.ErrEmailNotAvailable,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id"}).AddRow("user-xxxxx")
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE email = $1")).
					WithArgs("gopher@go.dev").
					WillReturnRows(mockRow)
			},
		},
		{
			name: "it should return error nil when email is not exist in database",
			args: args{
				ctx:   context.Background(),
				email: "gopher@go.dev",
			},
			expected: expected{
				err: nil,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE email = $1")).
					WithArgs("gopher@go.dev").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
			}

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, nil)
			err = repository.VerifyAvailableEmail(t.args.ctx, t.args.email)

			s.Equal(t.expected.err, err)
		})
	}
}

func (s *UserRepositoryTestSuite) TestFindByEmail() {
	type args struct {
		ctx   context.Context
		email string
	}
	type expected struct {
		user entity.User
		err  error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when database fail to query",
			args: args{
				ctx:   context.Background(),
				email: "gopher@go.dev",
			},
			expected: expected{
				user: entity.User{},
				err:  test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")).
					WithArgs("gopher@go.dev").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error ErrUserNotExist when user is not exist",
			args: args{
				ctx:   context.Background(),
				email: "gopher@go.dev",
			},
			expected: expected{
				user: entity.User{},
				err:  domain.ErrUserNotFound,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")).
					WithArgs("gopher@go.dev").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "it should return error nil and user when user is exist",
			args: args{
				ctx:   context.Background(),
				email: "gopher@go.dev",
			},
			expected: expected{
				user: entity.User{
					ID:        "user-xxxxx",
					Name:      "Gopher",
					Email:     "gopher@go.dev",
					Password:  "secret_password",
					CreatedAt: test.TimeBeforeNow,
					UpdatedAt: test.TimeBeforeNow,
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow("user-xxxxx", "Gopher", "gopher@go.dev", "secret_password", test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")).
					WithArgs("gopher@go.dev").
					WillReturnRows(mockRow)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
			}

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, nil)
			user, err := repository.FindByEmail(t.args.ctx, t.args.email)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.user, user)
		})
	}
}

func (s *UserRepositoryTestSuite) TestFindByID() {
	type args struct {
		ctx    context.Context
		userID entity.UserID
	}
	type expected struct {
		user entity.User
		err  error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when database fail to query",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				user: entity.User{},
				err:  test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1")).
					WithArgs("user-xxxxx").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error ErrUserNotFound when user is not found",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				user: entity.User{},
				err:  domain.ErrUserNotFound,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1")).
					WithArgs("user-xxxxx").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "it should return error nil and user when user is found",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				user: entity.User{
					ID:        "user-xxxxx",
					Name:      "Gopher",
					Email:     "gopher@go.dev",
					Password:  "secret_password",
					CreatedAt: test.TimeBeforeNow,
					UpdatedAt: test.TimeBeforeNow,
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow("user-xxxxx", "Gopher", "gopher@go.dev", "secret_password", test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1")).
					WithArgs("user-xxxxx").
					WillReturnRows(mockRow)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
			}

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, nil)
			user, err := repository.FindByID(t.args.ctx, t.args.userID)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.user, user)
		})
	}
}
