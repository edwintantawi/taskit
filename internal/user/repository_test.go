package user

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type RepositoryTestSuite struct {
	suite.Suite
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestCreate() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			ID:       entity.UserID("xxxxx"),
			Name:     "Gopher",
			Email:    "gopher@go.dev",
			Password: "secret_password",
		}

		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`)).
			WithArgs(u.ID, u.Name, u.Email, u.Password).
			WillReturnError(errors.New("database error"))

		mockUUID := &mocks.IDProvider{}
		mockUUID.On("Generate").Return(string(u.ID))

		repo := NewRepository(db, mockUUID)
		r, err := repo.Store(context.Background(), &u)

		s.Equal(errors.New("database error"), err)
		s.Empty(r)
	})

	s.Run("it should return user id when success insert to database", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			ID:       entity.UserID("xxxxx"),
			Name:     "Gopher",
			Email:    "gopher@go.dev",
			Password: "secret_password",
		}

		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`)).
			WithArgs(u.ID, u.Name, u.Email, u.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockUUID := &mocks.IDProvider{}
		mockUUID.On("Generate").Return(string(u.ID))

		repo := NewRepository(db, mockUUID)
		userID, err := repo.Store(context.Background(), &u)

		s.NoError(err)
		s.Equal(u.ID, userID)
	})
}

func (s *RepositoryTestSuite) TestVerifyAvailableEmail() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			Email: "gopher@go.dev",
		}

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE email = $1")).
			WithArgs(u.Email).
			WillReturnError(errors.New("database error"))

		repo := NewRepository(db, nil)
		err = repo.VerifyAvailableEmail(context.Background(), u.Email)

		s.Equal(errors.New("database error"), err)
	})

	s.Run("it should return error when email is not available", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			ID:    entity.UserID("xxxxx"),
			Email: "gopher@go.dev",
		}

		mockRow := sqlmock.NewRows([]string{"id"}).AddRow(u.ID)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE email = $1")).
			WithArgs(u.Email).
			WillReturnRows(mockRow)

		repo := NewRepository(db, nil)
		err = repo.VerifyAvailableEmail(context.Background(), u.Email)

		s.Equal(ErrEmailNotAvailable, err)
	})

	s.Run("it should return error nil when email is available", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			ID:    entity.UserID("xxxxx"),
			Email: "gopher@go.dev",
		}

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE email = $1")).
			WithArgs(u.Email).
			WillReturnError(sql.ErrNoRows)

		repo := NewRepository(db, nil)
		err = repo.VerifyAvailableEmail(context.Background(), u.Email)

		s.NoError(err)
	})
}
