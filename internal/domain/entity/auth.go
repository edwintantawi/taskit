package entity

import (
	"context"
	"errors"
	"time"
)

// Auth entity errors.
var (
	ErrAuthTokenExpired = errors.New("auth.entity.token.expired")
)

type AuthID string
type authUserIDKey string

// AuthUserIDKey is the key for the user_id value in the context.
const AuthUserIDKey = authUserIDKey("user_id")

// Auth represents an authentication in the system.
type Auth struct {
	ID        AuthID
	UserID    UserID
	Token     string
	ExpiresAt time.Time
}

// VerifyTokenExpires checks if the token has expired.
func (a *Auth) VerifyTokenExpires() error {
	if a.ExpiresAt.Before(time.Now()) {
		return ErrAuthTokenExpired
	}
	return nil
}

// GetAuthContext get the AuthUserIDKey from the context.
func GetAuthContext(ctx context.Context) UserID {
	userID := ctx.Value(AuthUserIDKey)
	if userID == nil {
		panic("Auth Context: Cannot get auth context, required context value user_id from auth middleware")
	}
	return userID.(UserID)
}
