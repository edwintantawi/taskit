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
	case domain.ErrEmailNotAvailable:
		return http.StatusBadRequest, "Email is not available"
	case domain.ErrUserEmailNotFound:
		return http.StatusBadRequest, "User email not found"
	case domain.ErrPasswordIncorrect:
		return http.StatusBadRequest, "Password is incorrect"
	default:
		return http.StatusInternalServerError, "Something went wrong"
	}
}
