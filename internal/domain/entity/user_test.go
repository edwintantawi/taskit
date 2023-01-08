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
	tests := []struct {
		name     string
		input    User
		expected error
	}{
		{name: "it should return error when email is invalid", input: User{Email: "invalid"}, expected: ErrEmailInvalid},
		{name: "it should return error when password is too short", input: User{Email: "gopher@go.dev", Password: "123"}, expected: ErrPasswordTooShort},
		{name: "it should return nil when all fields are valid", input: User{Email: "gopher@go.dev", Password: "123456", Name: "Gopher"}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}
