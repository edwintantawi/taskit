package domain

import "github.com/edwintantawi/taskit/internal/domain/entity"

// CreateUserIn represents the input of user creation.
type CreateUserIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserOut represents the output of user creation.
type CreateUserOut struct {
	ID    entity.UserID `json:"id"`
	Email string        `json:"email"`
}

// LoginAuthIn represent login input.
type LoginAuthIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginAuthOut represent login output.
type LoginAuthOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LogoutAuthIn represent logout input.
type LogoutAuthIn struct {
	RefreshToken string `json:"refresh_token"`
}

// GetProfileAuthIn represent profile input
type GetProfileAuthIn struct {
	UserID entity.UserID `json:"user_id"`
}

// GetProfileAuthOut represent profile output
type GetProfileAuthOut struct {
	ID    entity.UserID `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
}
