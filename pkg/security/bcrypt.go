package security

import "golang.org/x/crypto/bcrypt"

type bcryptx struct{}

func NewBcrypt() *bcryptx {
	return &bcryptx{}
}

// Hash hashes a raw string.
func (b *bcryptx) Hash(raw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
}
