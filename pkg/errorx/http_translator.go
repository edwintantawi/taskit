package errorx

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
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
	case entity.ErrEmailEmpty:
		return http.StatusBadRequest, "Email is required field"
	case entity.ErrEmailInvalid:
		return http.StatusBadRequest, "Email must be a valid email address"
	case entity.ErrPasswordEmpty:
		return http.StatusBadRequest, "Password is required field"
	case entity.ErrPasswordTooShort:
		return http.StatusBadRequest, fmt.Sprintf("Password must be greater then %d character in length", entity.MinPasswordLength)
	case entity.ErrNameEmpty:
		return http.StatusBadRequest, "Name is required field"
	// User repository
	case domain.ErrEmailNotAvailable:
		return http.StatusBadRequest, "Email is not available"
	case domain.ErrUserNotFound:
		return http.StatusNotFound, "User not found"
	// Auth entity
	case entity.ErrAuthTokenEmpty:
		return http.StatusBadRequest, "Refresh token is required field"
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
	// Task entity
	case entity.ErrContentEmpty:
		return http.StatusBadRequest, "Content is required field"
	// Task repository
	case domain.ErrTaskNotFound:
		return http.StatusNotFound, "Task not found"
	// Other
	default:
		return http.StatusInternalServerError, InternalServerErrorMessage
	}
}
