package idgen

import "github.com/google/uuid"

type uuidx struct{}

func NewUUID() *uuidx {
	return &uuidx{}
}

// Generate generates a new UUID string.
func (p *uuidx) Generate() string {
	return uuid.NewString()
}
