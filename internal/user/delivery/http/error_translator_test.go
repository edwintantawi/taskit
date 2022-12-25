package http

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/user/domain/entity"
	"github.com/edwintantawi/taskit/internal/user/repository"
)

type UserErrorTranslatorTestSuite struct {
	suite.Suite
}

func TestUserErrorTranslatorSuite(t *testing.T) {
	suite.Run(t, new(UserErrorTranslatorTestSuite))
}

func (s *UserErrorTranslatorTestSuite) TestErrorTranslator() {
	tests := []struct {
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{entity.ErrEmailEmpty, 400, "Email is required field"},
		{entity.ErrEmailInvalid, 400, "Email must be a valid email address"},
		{entity.ErrPasswordEmpty, 400, "Password is required field"},
		{entity.ErrPasswordTooShort, 400, fmt.Sprintf("Password must be greater then %d character in length", entity.MinPasswordLength)},
		{entity.ErrNameEmpty, 400, "Name is required field"},
		{repository.ErrEmailNotAvailable, 400, "Email is not available"},
		{errors.New("other error"), 500, "Something went wrong"},
	}

	for _, test := range tests {
		code, msg := ErrorTranslator(test.err)
		s.Equal(test.expectedCode, code)
		s.Equal(test.expectedMsg, msg)
	}
}
