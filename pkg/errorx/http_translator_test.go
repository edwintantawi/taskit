package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/security"
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
		{entity.ErrEmailInvalid, 400, "Email must be a valid email address"},
		{entity.ErrPasswordTooShort, 400, fmt.Sprintf("Password must be greater then %d character in length", entity.MinPasswordLength)},
		// User repository
		{domain.ErrEmailNotAvailable, 400, "Email is not available"},
		{domain.ErrUserNotFound, 404, "User not found"},
		// Auth entity
		{entity.ErrAuthTokenExpired, 400, "Refresh token is expired"},
		// Auth repository
		{domain.ErrAuthNotFound, 404, "Authentication not found"},
		// Auth usecase
		{domain.ErrPasswordIncorrect, 400, "Password is incorrect"},
		{domain.ErrEmailNotExist, 400, "Email is not exist"},
		// Task repository
		{domain.ErrTaskNotFound, 404, "Task not found"},
		// Task usecase
		{domain.ErrTaskAuthorization, 403, "Not have access to this task"},
		// DTO
		{dto.ErrEmailEmpty, 400, "Email is required field"},
		{dto.ErrPasswordEmpty, 400, "Password is required field"},
		{dto.ErrNameEmpty, 400, "Name is required field"},
		{dto.ErrRefreshTokenEmpty, 400, "Refresh token is required field"},
		{dto.ErrContentEmpty, 400, "Content is required field"},
		{dto.ErrTitleEmpty, 400, "Title is required field"},
		// Security JWT
		{security.ErrAccessTokenExpired, 401, "Access token is expired"},
		{security.ErrAccessTokenInvalid, 401, "Access token is invalid"},
		// Other
		{errors.New("other error"), 500, "Something went wrong"},
	}

	for _, test := range tests {
		code, msg := HTTPErrorTranslator(test.err)
		s.Equal(test.expectedCode, code)
		s.Equal(test.expectedMsg, msg)
	}
}
