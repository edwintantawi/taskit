package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrAuthTokenEmpty = errors.New("auth.entity.token.empty")
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

func (a *Auth) Validate() error {
	// remove all leading and trailing spaces
	a.Token = strings.TrimSpace(a.Token)

	switch {
	case a.Token == "":
		return ErrAuthTokenEmpty
	}
	return nil
}
