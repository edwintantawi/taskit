package dto

import "github.com/edwintantawi/taskit/internal/domain/entity"

// AuthLoginIn represent login input.
type AuthLoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthLoginIn) Validate() error {
	switch {
	case a.Email == "":
		return ErrEmailEmpty
	case a.Password == "":
		return ErrPasswordEmpty
	}
	return nil
}

// AuthLoginOut represent login output.
type AuthLoginOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthLogoutIn represent logout input.
type AuthLogoutIn struct {
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthLogoutIn) Validate() error {
	if a.RefreshToken == "" {
		return ErrRefreshTokenEmpty
	}
	return nil
}

// AuthProfileIn represent get profile input.
type AuthProfileIn struct {
	UserID entity.UserID `json:"-"`
}

// AuthProfileOut represent get profile output.
type AuthProfileOut struct {
	ID    entity.UserID `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
}

// AuthRefreshIn represent refresh input.
type AuthRefreshIn struct {
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthRefreshIn) Validate() error {
	if a.RefreshToken == "" {
		return ErrRefreshTokenEmpty
	}
	return nil
}

// AuthRefreshOut represent refresh output.
type AuthRefreshOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
