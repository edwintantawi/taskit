package security

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct{}

func NewBcrypt() Bcrypt {
	return Bcrypt{}
}

// Hash hashes a raw string.
func (b *Bcrypt) Hash(raw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
}

// Compare compares a raw string with a hashed string.
func (b *Bcrypt) Compare(raw string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}
