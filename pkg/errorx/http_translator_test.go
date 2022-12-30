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
		// User entity
		{entity.ErrEmailEmpty, 400, "Email is required field"},
		{entity.ErrEmailInvalid, 400, "Email must be a valid email address"},
		{entity.ErrPasswordEmpty, 400, "Password is required field"},
		{entity.ErrPasswordTooShort, 400, fmt.Sprintf("Password must be greater then %d character in length", entity.MinPasswordLength)},
		{entity.ErrNameEmpty, 400, "Name is required field"},
		// User repository
		{domain.ErrEmailNotAvailable, 400, "Email is not available"},
		{domain.ErrUserNotFound, 404, "User not found"},
		// Auth entity
		{entity.ErrAuthTokenEmpty, 400, "Refresh token is required field"},
		{entity.ErrAuthTokenExpired, 400, "Refresh token is expired"},
		// Auth repository
		{domain.ErrAuthNotFound, 404, "Authentication not found"},
		// Auth usecase
		{domain.ErrPasswordIncorrect, 400, "Password is incorrect"},
		{domain.ErrEmailNotExist, 400, "Email is not exist"},
		// Task entity
		{entity.ErrContentEmpty, 400, "Content is required field"},
		// Task repository
		{domain.ErrTaskNotFound, 404, "Task not found"},
		// Task usecase
		{domain.ErrTaskAuthorization, 403, "Not have access to this task"},
		// Other
		{errors.New("other error"), 500, "Something went wrong"},
	}

	for _, test := range tests {
		code, msg := HTTPErrorTranslator(test.err)
		s.Equal(test.expectedCode, code)
		s.Equal(test.expectedMsg, msg)
	}
}
