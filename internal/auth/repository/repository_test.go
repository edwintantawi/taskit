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

type AuthRepositoryTestSuite struct {
	suite.Suite
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) TestStore() {
	s.Run("it should return error when database failed to store token", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		a := entity.Auth{
			ID:     "yyyyy",
			UserID: "xxxxx",
			Token:  "xxxxx.xxxxx.xxxxx",
		}

		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO authentications (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`)).
			WithArgs(a.ID, a.UserID, a.Token, a.ExpiresAt).
			WillReturnError(errors.New("database error"))

		mockIDProvider := &mocks.IDProvider{}
		mockIDProvider.On("Generate").Return(string(a.ID))

		repo := New(db, mockIDProvider)
		err = repo.Store(context.Background(), &a)

		s.Equal(err, errors.New("database error"))
	})

	s.Run("it should return error nil when successfully store token", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		a := entity.Auth{
			ID:     "yyyyy",
			UserID: "xxxxx",
			Token:  "xxxxx.xxxxx.xxxxx",
		}

		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO authentications (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`)).
			WithArgs(a.ID, a.UserID, a.Token, a.ExpiresAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockIDProvider := &mocks.IDProvider{}
		mockIDProvider.On("Generate").Return(string(a.ID))

		repo := New(db, mockIDProvider)
		err = repo.Store(context.Background(), &a)

		s.NoError(err)
	})
}

func (s *AuthRepositoryTestSuite) TestDelete() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		a := entity.Auth{UserID: "xxxxx", Token: "xxxxx.xxxxx.xxxxx"}
		mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM authentications WHERE token = $1`)).
			WithArgs(a.Token).
			WillReturnError(errors.New("database error"))

		repo := New(db, nil)
		err = repo.Delete(context.Background(), &a)

		s.Equal(errors.New("database error"), err)
	})

	s.Run("it should return error when user id or token not found", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		a := entity.Auth{UserID: "xxxxx", Token: "xxxxx.xxxxx.xxxxx"}
		mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM authentications WHERE token = $1`)).
			WithArgs(a.Token).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := New(db, nil)
		err = repo.Delete(context.Background(), &a)

		s.Equal(domain.ErrAuthNotExist, err)
	})

	s.Run("it should return error when fail to get rows affected", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		a := entity.Auth{UserID: "xxxxx", Token: "xxxxx.xxxxx.xxxxx"}
		mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM authentications WHERE token = $1`)).
			WithArgs(a.Token).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))

		repo := New(db, nil)
		err = repo.Delete(context.Background(), &a)

		s.Equal(errors.New("unexpected error"), err)
	})

	s.Run("it should return error nil when authentication deleted successfully", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		a := entity.Auth{UserID: "xxxxx", Token: "xxxxx.xxxxx.xxxxx"}

		mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM authentications WHERE token = $1`)).
			WithArgs(a.Token).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(db, nil)
		err = repo.Delete(context.Background(), &a)

		s.NoError(err)
	})
}

func (s *AuthRepositoryTestSuite) TestFindByToken() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		token := "xxxxx.xxxxx.xxxxx"

		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
			WithArgs(token).
			WillReturnError(errors.New("database error"))

		repo := New(db, nil)
		_, err = repo.FindByToken(context.Background(), token)

		s.Equal(errors.New("database error"), err)
	})

	s.Run("it should return error when token not found", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		token := "xxxxx.xxxxx.xxxxx"

		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
			WithArgs(token).
			WillReturnError(sql.ErrNoRows)

		repo := New(db, nil)
		_, err = repo.FindByToken(context.Background(), token)

		s.Equal(domain.ErrAuthNotExist, err)
	})

	s.Run("it should return error when fail to scan row", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		token := "xxxxx.xxxxx.xxxxx"

		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
			WithArgs(token).
			WillReturnError(errors.New("unexpected error"))

		repo := New(db, nil)
		_, err = repo.FindByToken(context.Background(), token)

		s.Equal(errors.New("unexpected error"), err)
	})

	s.Run("it should return error nil when authentication found", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		token := "xxxxx.xxxxx.xxxxx"

		mockRow := sqlmock.NewRows([]string{"id", "user_id", "token", "expires_at"}).AddRow("xxxxx", "xxxxx", "xxxxx.xxxxx.xxxxx", time.Now().Add(time.Hour))
		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`)).
			WithArgs(token).
			WillReturnRows(mockRow)

		repo := New(db, nil)
		auth, err := repo.FindByToken(context.Background(), token)

		s.NoError(err)
		s.NotEmpty(auth)
	})
}
