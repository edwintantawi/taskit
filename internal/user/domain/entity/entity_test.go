package entity

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserEntityTestSuite struct {
	suite.Suite
}

func TestUserEntitySuite(t *testing.T) {
	suite.Run(t, new(UserEntityTestSuite))
}

func (s *UserEntityTestSuite) TestValidate() {
	s.Run("it should return error when email is empty", func() {
		u := User{}
		err := u.Validate()

		s.Equal(ErrEmailEmpty, err)
	})

	s.Run("it should return error when email is invalid", func() {
		u := User{
			Email: "invalid",
		}
		err := u.Validate()

		s.Equal(ErrEmailInvalid, err)
	})

	s.Run("it should return error when password is empty", func() {
		u := User{
			Email: "gopher@go.dev",
		}
		err := u.Validate()

		s.Equal(ErrPasswordEmpty, err)
	})

	s.Run("it should return error when password is too short", func() {
		u := User{
			Email:    "gopher@go.dev",
			Password: "123",
		}
		err := u.Validate()

		s.Equal(ErrPasswordTooShort, err)
	})

	s.Run("it should return error when name is empty", func() {
		u := User{
			Email:    "gopher@go.dev",
			Password: "123456",
		}
		err := u.Validate()

		s.Equal(ErrNameEmpty, err)
	})

	s.Run("it should return nil when all fields are valid", func() {
		u := User{
			Email:    "gopher@go.dev",
			Password: "123456",
			Name:     "Gopher",
		}
		err := u.Validate()

		s.Nil(err)
	})
}
