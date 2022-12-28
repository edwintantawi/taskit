package entity

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrAuthTokenEmpty   = errors.New("auth.entity.token.empty")
	ErrAuthTokenExpired = errors.New("auth.entity.token.expired")
)

type AuthID string
type authUserIDKey string

const AuthUserIDKey = authUserIDKey("user_id")

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

func GetAuthContext(ctx context.Context) UserID {
	userID := ctx.Value(AuthUserIDKey)
	if userID == nil {
		panic("Auth Context: Cannot get auth context, required context value user_id from auth middleware")
	}
	return userID.(UserID)
}
