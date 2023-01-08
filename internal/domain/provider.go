package domain

import (
	"time"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// IDProvider represent id generator contract
type IDProvider interface {
	Generate() string
}

// HashProvider represent hasher contract
type HashProvider interface {
	Hash(raw string) ([]byte, error)
	Compare(raw string, hashed string) error
}

// JWTProvider represent jwt generator contract.
type JWTProvider interface {
	GenerateAccessToken(userID entity.UserID) (string, time.Time, error)
	GenerateRefreshToken(userID entity.UserID) (string, time.Time, error)
	VerifyAccessToken(rawToken string) (entity.UserID, error)
}

// Validater represent object with validate method.
type Validater interface {
	Validate() error
}

// ValidatorProvider represent validator contract.
type ValidatorProvider interface {
	Validate(validater Validater) error
}
