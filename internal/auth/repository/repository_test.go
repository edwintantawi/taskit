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

type AuthRepositoryTestSuite struct {
	suite.Suite
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

type dependency struct {
	mockDB     sqlmock.Sqlmock
	idProvider *mocks.IDProvider
}

func (s *AuthRepositoryTestSuite) TestStore() {
	type args struct {
		ctx  context.Context
		auth *entity.Auth
	}
	tests := []struct {
		name     string
		args     args
		expected error
		setup    func(d *dependency)
	}{
		{
			name:     "it should return error when database fail to store",
			args:     args{ctx: context.Background(), auth: &entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}},
			expected: test.ErrDatabase,
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return(string("auth-xxxxx"))
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO authentications (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`)).
					WithArgs("auth-xxxxx", "user-xxxxx", "yyyyy.yyyyy.yyyyy", test.TimeAfterNow).
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name:     "it should return error nil when successfully store",
			args:     args{ctx: context.Background(), auth: &entity.Auth{UserID: "user-xxxxx", Token: "yyyyy.yyyyy.yyyyy", ExpiresAt: test.TimeAfterNow}},
			expected: nil,
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return(string("auth-xxxxx"))
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO authentications (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`)).
					WithArgs("auth-xxxxx", "user-xxxxx", "yyyyy.yyyyy.yyyyy", test.TimeAfterNow).
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
			err = repository.Store(t.args.ctx, t.args.auth)

			s.Equal(t.expected, err)
		})
	}
}

func (s *AuthRepositoryTestSuite) TestVerifyAvailableByID() {
	type args struct {
		ctx   context.Context
		token string
	}
	type expected struct {
		error error
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
				token: "yyyyy.yyyyy.yyyyy",
			},
			expected: expected{
				error: test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error when authentication not found",
			args: args{
				ctx:   context.Background(),
				token: "yyyyy.yyyyy.yyyyy",
			},
			expected: expected{
				error: domain.ErrAuthNotFound,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "it should return error when database fail to scan",
			args: args{
				ctx:   context.Background(),
				token: "yyyyy.yyyyy.yyyyy",
			},
			expected: expected{
				error: test.ErrRowScan,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(test.ErrRowScan)
			},
		},
		{
			name: "it should return error nil when authentication found",
			args: args{
				ctx:   context.Background(),
				token: "yyyyy.yyyyy.yyyyy",
			},
			expected: expected{
				error: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id"}).AddRow("auth-xxxxx")
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
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
			err = repository.VerifyAvailableByToken(t.args.ctx, t.args.token)

			s.Equal(t.expected.error, err)
		})
	}
}

func (s *AuthRepositoryTestSuite) TestDelete() {
	type args struct {
		ctx  context.Context
		auth *entity.Auth
	}
	tests := []struct {
		name     string
		args     args
		expected error
		setup    func(d *dependency)
	}{
		{
			name:     "it should return error when database fail to delete",
			args:     args{ctx: context.Background(), auth: &entity.Auth{Token: "yyyyy.yyyyy.yyyyy"}},
			expected: test.ErrDatabase,
			setup: func(d *dependency) {
				d.mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name:     "it should return error nil when successfully delete",
			args:     args{ctx: context.Background(), auth: &entity.Auth{Token: "yyyyy.yyyyy.yyyyy"}},
			expected: nil,
			setup: func(d *dependency) {
				d.mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
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

			d := &dependency{
				mockDB: mockDB,
			}
			t.setup(d)

			repository := New(db, nil)
			err = repository.Delete(t.args.ctx, t.args.auth)

			s.Equal(t.expected, err)
		})
	}
}

func (s *AuthRepositoryTestSuite) TestFindByToken() {
	type args struct {
		ctx   context.Context
		token string
	}
	type expected struct {
		auth entity.Auth
		err  error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when database fail to find",
			args: args{ctx: context.Background(), token: "yyyyy.yyyyy.yyyyy"},
			expected: expected{
				auth: entity.Auth{},
				err:  test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error ErrAuthNotFound when row not found",
			args: args{ctx: context.Background(), token: "yyyyy.yyyyy.yyyyy"},
			expected: expected{
				auth: entity.Auth{},
				err:  domain.ErrAuthNotFound,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "it should return error when fail to scan row",
			args: args{ctx: context.Background(), token: "yyyyy.yyyyy.yyyyy"},
			expected: expected{
				auth: entity.Auth{},
				err:  test.ErrRowScan,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
					WillReturnError(test.ErrRowScan)
			},
		},
		{
			name: "it should return error nil and authentication when found",
			args: args{ctx: context.Background(), token: "yyyyy.yyyyy.yyyyy"},
			expected: expected{
				auth: entity.Auth{
					ID:        "auth-xxxxx",
					UserID:    "user-xxxxx",
					Token:     "yyyyy.yyyyy.yyyyy",
					ExpiresAt: test.TimeAfterNow,
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "user_id", "token", "expires_at"}).
					AddRow("auth-xxxxx", "user-xxxxx", "yyyyy.yyyyy.yyyyy", test.TimeAfterNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
					WithArgs("yyyyy.yyyyy.yyyyy").
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

			d := &dependency{
				mockDB: mockDB,
			}
			t.setup(d)

			repository := New(db, nil)
			auth, err := repository.FindByToken(t.args.ctx, t.args.token)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.auth, auth)
		})
	}
}
