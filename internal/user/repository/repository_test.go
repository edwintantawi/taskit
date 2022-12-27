package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type UserRepositoryTestSuite struct {
	suite.Suite
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestCreate() {
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

		mockIDProvider := &mocks.IDProvider{}
		mockIDProvider.On("Generate").Return(string(u.ID))

		repo := New(db, mockIDProvider)
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

		mockIDProvider := &mocks.IDProvider{}
		mockIDProvider.On("Generate").Return(string(u.ID))

		repo := New(db, mockIDProvider)
		userID, err := repo.Store(context.Background(), &u)

		s.NoError(err)
		s.Equal(u.ID, userID)
	})
}

func (s *UserRepositoryTestSuite) TestVerifyAvailableEmail() {
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

		repo := New(db, nil)
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

		repo := New(db, nil)
		err = repo.VerifyAvailableEmail(context.Background(), u.Email)

		s.Equal(domain.ErrEmailNotAvailable, err)
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

		repo := New(db, nil)
		err = repo.VerifyAvailableEmail(context.Background(), u.Email)

		s.NoError(err)
	})
}

func (s *UserRepositoryTestSuite) TestFindByEmail() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			Email: "gopher@go.dev",
		}

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")).
			WithArgs(u.Email).
			WillReturnError(errors.New("database error"))

		repo := New(db, nil)
		user, err := repo.FindByEmail(context.Background(), u.Email)

		s.Equal(errors.New("database error"), err)
		s.Empty(user)
	})

	s.Run("it should return error user not found when user is not exist", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			Email: "gopher@go.dev",
		}

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")).
			WithArgs(u.Email).
			WillReturnError(sql.ErrNoRows)

		repo := New(db, nil)
		user, err := repo.FindByEmail(context.Background(), u.Email)

		s.Equal(domain.ErrUserEmailNotExist, err)
		s.Empty(user)
	})

	s.Run("it should return user when user is exist", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			ID:        "xxxxx",
			Name:      "Gopher",
			Email:     "gopher@go.dev",
			Password:  "secret_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRow := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).AddRow(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1")).
			WithArgs(u.Email).
			WillReturnRows(mockRow)

		repo := New(db, nil)
		user, err := repo.FindByEmail(context.Background(), u.Email)

		s.NoError(err)
		s.Equal(u, user)
	})
}

func (s *UserRepositoryTestSuite) TestFindByID() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		userID := "xxxxx"

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1")).
			WithArgs(userID).
			WillReturnError(errors.New("database error"))

		repo := New(db, nil)
		user, err := repo.FindByID(context.Background(), entity.UserID(userID))

		s.Equal(errors.New("database error"), err)
		s.Empty(user)
	})

	s.Run("it should return error user not found when user is not exist", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		userID := "xxxxx"

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1")).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		repo := New(db, nil)
		user, err := repo.FindByID(context.Background(), entity.UserID(userID))

		s.Equal(domain.ErrUserIDNotExist, err)
		s.Empty(user)
	})

	s.Run("it should return user when user is exist", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		u := entity.User{
			ID:        "xxxxx",
			Name:      "Gopher",
			Email:     "gopher@go.dev",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRow := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).AddRow(u.ID, u.Name, u.Email, u.CreatedAt, u.UpdatedAt)
		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1")).
			WithArgs(u.ID).
			WillReturnRows(mockRow)

		repo := New(db, nil)
		user, err := repo.FindByID(context.Background(), entity.UserID(u.ID))

		s.NoError(err)
		s.Equal(u.ID, user.ID)
		s.Equal(u.Name, user.Name)
		s.Equal(u.Email, user.Email)
		s.Equal(u.CreatedAt, user.CreatedAt)
		s.Equal(u.UpdatedAt, user.UpdatedAt)
	})
}
