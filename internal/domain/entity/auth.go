package entity

import "time"

type AuthID string
type AuthUserIDKey string

// Auth represents an authentication in the system.
type Auth struct {
	ID        AuthID
	UserID    UserID
	Token     string
	ExpiresAt time.Time
}
