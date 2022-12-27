package usecase

import (
	"context"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type usecase struct {
	authRepository domain.AuthRepository
	userRepository domain.UserRepository
	hashProvider   domain.HashProvider
	jwtProvider    domain.JWTProvider
}

// New create a new auth usecase.
func New(authRepository domain.AuthRepository, userRepository domain.UserRepository, hashProvider domain.HashProvider, jwtProvider domain.JWTProvider) domain.AuthUsecase {
	return &usecase{
		authRepository: authRepository,
		userRepository: userRepository,
		hashProvider:   hashProvider,
		jwtProvider:    jwtProvider,
	}
}

// Login authenticates a user.
func (u *usecase) Login(ctx context.Context, payload *domain.LoginAuthIn) (domain.LoginAuthOut, error) {
	user, err := u.userRepository.FindByEmail(ctx, payload.Email)
	if err != nil {
		return domain.LoginAuthOut{}, err
	}

	if err := u.hashProvider.Compare(payload.Password, user.Password); err != nil {
		return domain.LoginAuthOut{}, domain.ErrPasswordIncorrect
	}

	accessToken, _, err := u.jwtProvider.GenerateAccessToken(user.ID)
	if err != nil {
		return domain.LoginAuthOut{}, err
	}
	refreshToken, expires, err := u.jwtProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return domain.LoginAuthOut{}, err
	}

	auth := &entity.Auth{UserID: user.ID, Token: refreshToken, ExpiresAt: expires}
	if err := u.authRepository.Store(ctx, auth); err != nil {
		return domain.LoginAuthOut{}, err
	}

	return domain.LoginAuthOut{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Logout remove user authentication.
func (u *usecase) Logout(ctx context.Context, payload *domain.LogoutAuthIn) error {
	auth := &entity.Auth{Token: payload.RefreshToken}
	if err := auth.Validate(); err != nil {
		return err
	}

	err := u.authRepository.Delete(ctx, auth)
	if err != nil {
		return err
	}

	return nil
}

// GetProfile get user authenticated profile.
func (u *usecase) GetProfile(ctx context.Context, payload *domain.GetProfileAuthIn) (domain.GetProfileAuthOut, error) {
	user, err := u.userRepository.FindByID(ctx, payload.UserID)
	if err != nil {
		return domain.GetProfileAuthOut{}, err
	}
	return domain.GetProfileAuthOut{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

// Refresh refresh user authentication token.
func (u *usecase) Refresh(ctx context.Context, payload *domain.RefreshAuthIn) (domain.RefreshAuthOut, error) {
	auth, err := u.authRepository.FindByToken(ctx, payload.RefreshToken)
	if err != nil {
		return domain.RefreshAuthOut{}, err
	}
	if err := auth.VerifyTokenExpires(); err != nil {
		return domain.RefreshAuthOut{}, err
	}

	accessToken, _, err := u.jwtProvider.GenerateAccessToken(auth.UserID)
	if err != nil {
		return domain.RefreshAuthOut{}, err
	}
	refreshToken, expires, err := u.jwtProvider.GenerateRefreshToken(auth.UserID)
	if err != nil {
		return domain.RefreshAuthOut{}, err
	}

	if err := u.authRepository.Delete(ctx, &auth); err != nil {
		return domain.RefreshAuthOut{}, err
	}

	newAuth := &entity.Auth{UserID: auth.UserID, Token: refreshToken, ExpiresAt: expires}
	if err := u.authRepository.Store(ctx, newAuth); err != nil {
		return domain.RefreshAuthOut{}, err
	}

	return domain.RefreshAuthOut{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
