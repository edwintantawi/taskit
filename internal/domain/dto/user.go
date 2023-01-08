package dto

import (
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// CreateUserIn represents the input of user creation.
type CreateUserIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *CreateUserIn) Validate() error {
	switch {
	case d.Email == "":
		return ErrEmailEmpty
	case d.Password == "":
		return ErrPasswordEmpty
	case d.Name == "":
		return ErrNameEmpty
	}
	return nil
}

// CreateUserOut represents the output of user creation.
type CreateUserOut struct {
	ID    entity.UserID `json:"id"`
	Email string        `json:"email"`
}
