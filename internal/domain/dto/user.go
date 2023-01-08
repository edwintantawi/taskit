package dto

import (
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// UserCreateIn represents the input of user creation.
type UserCreateIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserCreateIn) Validate() error {
	switch {
	case u.Email == "":
		return ErrEmailEmpty
	case u.Password == "":
		return ErrPasswordEmpty
	case u.Name == "":
		return ErrNameEmpty
	}
	return nil
}

// UserCreateOut represents the output of user creation.
type UserCreateOut struct {
	ID    entity.UserID `json:"id"`
	Email string        `json:"email"`
}
