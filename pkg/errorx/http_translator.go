package errorx

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// HTTPErrorTranslator translates error to http status code and human readable error message.
func HTTPErrorTranslator(err error) (code int, msg string) {
	log.Println("[ERROR]", err)
	switch err {
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
	case entity.ErrAuthTokenEmpty:
		return http.StatusBadRequest, "Refresh token is required field"
	case entity.ErrAuthTokenExpired:
		return http.StatusBadRequest, "Refresh token is expired"
	case domain.ErrEmailNotAvailable:
		return http.StatusBadRequest, "Email is not available"
	case domain.ErrUserEmailNotExist:
		return http.StatusBadRequest, "User email not found"
	case domain.ErrPasswordIncorrect:
		return http.StatusBadRequest, "Password is incorrect"
	case domain.ErrAuthNotExist:
		return http.StatusBadRequest, "Authentication token not exist"
	case domain.ErrUserIDNotExist:
		return http.StatusNotFound, "User not found"
	default:
		return http.StatusInternalServerError, "Something went wrong"
	}
}
