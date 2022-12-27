package entity

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthEntityTestSuite struct {
	suite.Suite
}

func TestAuthEntitySuite(t *testing.T) {
	suite.Run(t, new(AuthEntityTestSuite))
}

func (s *AuthEntityTestSuite) TestValidate() {
	s.Run("it should return error when token is empty", func() {
		u := Auth{}
		err := u.Validate()

		s.Equal(ErrAuthTokenEmpty, err)
	})

	s.Run("it should return nil when all fields are valid", func() {
		u := Auth{
			Token: "xxxxx.xxxxx.xxxxx",
		}
		err := u.Validate()

		s.Nil(err)
	})
}
