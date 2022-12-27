package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type HTTPErrorTranslatorTestSuite struct {
	suite.Suite
}

func TestHTTPErrorTranslatorSuite(t *testing.T) {
	suite.Run(t, new(HTTPErrorTranslatorTestSuite))
}

func (s *HTTPErrorTranslatorTestSuite) TestErrorTranslator() {
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
		{entity.ErrAuthTokenEmpty, 400, "Refresh token is required field"},
		{domain.ErrEmailNotAvailable, 400, "Email is not available"},
		{domain.ErrUserEmailNotExist, 400, "User email not found"},
		{domain.ErrPasswordIncorrect, 400, "Password is incorrect"},
		{domain.ErrAuthNotExist, 400, "Authentication token not exist"},
		{errors.New("other error"), 500, "Something went wrong"},
	}

	for _, test := range tests {
		code, msg := HTTPErrorTranslator(test.err)
		s.Equal(test.expectedCode, code)
		s.Equal(test.expectedMsg, msg)
	}
}
