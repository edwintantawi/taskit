package domain

import "time"

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
	GenerateAccessToken(payload map[string]interface{}) (string, time.Time, error)
	GenerateRefreshToken(payload map[string]interface{}) (string, time.Time, error)
}
