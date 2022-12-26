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

// Compare compares a raw string with a hashed string.
func (b *bcryptx) Compare(raw string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}
