package idgen

import "github.com/google/uuid"

type UUID struct{}

// NewUUID creates a new UUID generator.
func NewUUID() UUID {
	return UUID{}
}

// Generate generates a new UUID string.
func (p *UUID) Generate() string {
	return uuid.NewString()
}
