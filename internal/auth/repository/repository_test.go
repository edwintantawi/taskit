package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) TestHere() {
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
