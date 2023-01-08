package dto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthDTOTestSuite struct {
	suite.Suite
}

func TestAuthDTOSuite(t *testing.T) {
	suite.Run(t, new(AuthDTOTestSuite))
}

func (s *AuthDTOTestSuite) TestAuthLoginIn() {
	tests := []struct {
		name     string
		input    AuthLoginIn
		expected error
	}{
		{name: "it should return error when email is empty", input: AuthLoginIn{}, expected: ErrEmailEmpty},
		{name: "it should return error when password is empty", input: AuthLoginIn{Email: "gopher@go.dev"}, expected: ErrPasswordEmpty},
		{name: "it should return nil when all fields are valid", input: AuthLoginIn{Email: "gopher@go.dev", Password: "secret_password"}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}

func (s *AuthDTOTestSuite) TestAuthLogoutIn() {
	tests := []struct {
		name     string
		input    AuthLogoutIn
		expected error
	}{
		{name: "it should return error when password is empty", input: AuthLogoutIn{RefreshToken: ""}, expected: ErrRefreshTokenEmpty},
		{name: "it should return nil when all fields are valid", input: AuthLogoutIn{RefreshToken: "yyyyy.yyyyy.yyyyy"}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}

func (s *AuthDTOTestSuite) TestAuthRefreshIn() {
	tests := []struct {
		name     string
		input    AuthRefreshIn
		expected error
	}{
		{name: "it should return error when password is empty", input: AuthRefreshIn{RefreshToken: ""}, expected: ErrRefreshTokenEmpty},
		{name: "it should return nil when all fields are valid", input: AuthRefreshIn{RefreshToken: "yyyyy.yyyyy.yyyyy"}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := test.input.Validate()
			s.Equal(test.expected, err)
		})
	}
}
