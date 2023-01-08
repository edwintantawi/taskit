package entity

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AuthEntityTestSuite struct {
	suite.Suite
}

func TestAuthEntitySuite(t *testing.T) {
	suite.Run(t, new(AuthEntityTestSuite))
}

func (s *AuthEntityTestSuite) TestVerifyTokenExpires() {
	tests := []struct {
		name     string
		input    Auth
		expected error
	}{
		{name: "it should return error when auth is expired", input: Auth{ExpiresAt: time.Now().Add(-1 * time.Hour)}, expected: ErrAuthTokenExpired},
		{name: "it should return nill when auth is not expired", input: Auth{ExpiresAt: time.Now().Add(1 * time.Hour)}, expected: nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.Equal(test.expected, test.input.VerifyTokenExpires())
		})
	}
}

func (s *AuthEntityTestSuite) TestGetAuthContext() {
	s.Run("it should panic when auth context is not set", func() {
		s.Panics(func() {
			GetAuthContext(context.Background())
		})
	})

	s.Run("it should return user id when auth context is set", func() {
		userID := UserID("xxxxx")
		ctx := context.WithValue(context.Background(), AuthUserIDKey, userID)
		s.Equal(userID, GetAuthContext(ctx))
	})
}
