package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrAuthTokenEmpty   = errors.New("auth.entity.token.empty")
	ErrAuthTokenExpired = errors.New("auth.entity.token.expired")
)

type AuthID string
type AuthUserIDKey string

// Auth represents an authentication in the system.
type Auth struct {
	ID        AuthID
	UserID    UserID
	Token     string
	ExpiresAt time.Time
}

// Validate auth fields.
func (a *Auth) Validate() error {
	// remove all leading and trailing spaces
	a.Token = strings.TrimSpace(a.Token)

	switch {
	case a.Token == "":
		return ErrAuthTokenEmpty
	}
	return nil
}

// VerifyTokenExpires checks if the token has expired.
func (a *Auth) VerifyTokenExpires() error {
	if a.ExpiresAt.Before(time.Now()) {
		return ErrAuthTokenExpired
	}
	return nil
}
