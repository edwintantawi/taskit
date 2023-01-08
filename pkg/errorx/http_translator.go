package errorx

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// HTTPError message
var (
	InternalServerErrorMessage = "Something went wrong"
)

// HTTPErrorTranslator translates error to http status code and human readable error message.
func HTTPErrorTranslator(err error) (code int, msg string) {
	log.Println("[ERROR]", err)
	switch err {
	// User entity
	case entity.ErrEmailInvalid:
		return http.StatusBadRequest, "Email must be a valid email address"
	case entity.ErrPasswordTooShort:
		return http.StatusBadRequest, fmt.Sprintf("Password must be greater then %d character in length", entity.MinPasswordLength)
	// User repository
	case domain.ErrEmailNotAvailable:
		return http.StatusBadRequest, "Email is not available"
	case domain.ErrUserNotFound:
		return http.StatusNotFound, "User not found"
	// Auth entity
	case entity.ErrAuthTokenExpired:
		return http.StatusBadRequest, "Refresh token is expired"
	// Auth repository
	case domain.ErrAuthNotFound:
		return http.StatusNotFound, "Authentication not found"
	case domain.ErrEmailNotExist:
		return http.StatusBadRequest, "Email is not exist"
	// Auth usecase
	case domain.ErrPasswordIncorrect:
		return http.StatusBadRequest, "Password is incorrect"
	// Task repository
	case domain.ErrTaskNotFound:
		return http.StatusNotFound, "Task not found"
	// Task usecase
	case domain.ErrTaskAuthorization:
		return http.StatusForbidden, "Not have access to this task"
	// DTO
	case dto.ErrEmailEmpty:
		return http.StatusBadRequest, "Email is required field"
	case dto.ErrPasswordEmpty:
		return http.StatusBadRequest, "Password is required field"
	case dto.ErrNameEmpty:
		return http.StatusBadRequest, "Name is required field"
	case dto.ErrRefreshTokenEmpty:
		return http.StatusBadRequest, "Refresh token is required field"
	case dto.ErrTaskContentEmpty:
		return http.StatusBadRequest, "Content is required field"
	// Other
	default:
		return http.StatusInternalServerError, InternalServerErrorMessage
	}
}
