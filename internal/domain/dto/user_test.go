package dto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserDTOTestSuite struct {
	suite.Suite
}

func TestUserDTOSuite(t *testing.T) {
	suite.Run(t, new(UserDTOTestSuite))
}

func (s *UserDTOTestSuite) TestUserCreateIn() {
	tests := []struct {
		name     string
		input    UserCreateIn
		expected error
	}{
		{name: "it should return error when email is empty", input: UserCreateIn{}, expected: ErrEmailEmpty},
		{name: "it should return error when password is empty", input: UserCreateIn{Email: "gopher@go.dev"}, expected: ErrPasswordEmpty},
		{name: "it should return error when name is empty", input: UserCreateIn{Email: "gopher@go.dev", Password: "123456"}, expected: ErrNameEmpty},
		{name: "it should return nil when all fields are valid", input: UserCreateIn{Email: "gopher@go.dev", Password: "123456", Name: "Gopher"}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}
