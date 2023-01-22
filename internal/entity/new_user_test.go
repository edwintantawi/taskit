package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser_Validate(t *testing.T) {
	tests := map[string]struct {
		args     NewUser
		expected error
	}{
		"it should return error ErrEmptyName when name is empty": {
			args:     NewUser{},
			expected: ErrEmptyName,
		},
		"it should return error ErrInvalidEmail when email is not valid": {
			args: NewUser{
				Name:  "gopher",
				Email: "invalid-email",
			},
			expected: ErrInvalidEmail,
		},
		"it should return error ErrTooShortPassword when the password less than 6 characters": {
			args: NewUser{
				Name:     "gopher",
				Email:    "gopher@go.dev",
				Password: "12345",
			},
			expected: ErrTooShortPassword,
		},
		"it should return error nil when no validation error": {
			args: NewUser{
				Name:     "gopher",
				Email:    "gopher@go.dev",
				Password: "123456",
			},
			expected: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.args.Validate()
			assert.Equal(t, tc.expected, err)
		})
	}
}
