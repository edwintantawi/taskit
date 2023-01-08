package usecase

import (
	"context"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type usecase struct {
	authRepository domain.AuthRepository
	userRepository domain.UserRepository
	hashProvider   domain.HashProvider
	jwtProvider    domain.JWTProvider
}

// New create a new auth usecase.
func New(
	authRepository domain.AuthRepository,
	userRepository domain.UserRepository,
	hashProvider domain.HashProvider,
	jwtProvider domain.JWTProvider,
) domain.AuthUsecase {
	return &usecase{
		authRepository: authRepository,
		userRepository: userRepository,
		hashProvider:   hashProvider,
		jwtProvider:    jwtProvider,
	}
}

// Login authenticates a user.
func (u *usecase) Login(ctx context.Context, payload *dto.AuthLoginIn) (dto.AuthLoginOut, error) {
	user, err := u.userRepository.FindByEmail(ctx, payload.Email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return dto.AuthLoginOut{}, domain.ErrEmailNotExist
	} else if err != nil {
		return dto.AuthLoginOut{}, err
	}

	if err := u.hashProvider.Compare(payload.Password, user.Password); err != nil {
		return dto.AuthLoginOut{}, domain.ErrPasswordIncorrect
	}

	accessToken, _, err := u.jwtProvider.GenerateAccessToken(user.ID)
	if err != nil {
		return dto.AuthLoginOut{}, err
	}
	refreshToken, expires, err := u.jwtProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthLoginOut{}, err
	}

	auth := &entity.Auth{UserID: user.ID, Token: refreshToken, ExpiresAt: expires}
	if err := u.authRepository.Store(ctx, auth); err != nil {
		return dto.AuthLoginOut{}, err
	}

	return dto.AuthLoginOut{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Logout remove user authentication.
func (u *usecase) Logout(ctx context.Context, payload *dto.AuthLogoutIn) error {
	auth := &entity.Auth{Token: payload.RefreshToken}

	if err := u.authRepository.VerifyAvailableByToken(ctx, auth.Token); err != nil {
		return err
	}
	if err := u.authRepository.DeleteByToken(ctx, auth.Token); err != nil {
		return err
	}

	return nil
}

// GetProfile get user authenticated profile.
func (u *usecase) GetProfile(ctx context.Context, payload *dto.AuthProfileIn) (dto.AuthProfileOut, error) {
	user, err := u.userRepository.FindByID(ctx, payload.UserID)
	if err != nil {
		return dto.AuthProfileOut{}, err
	}
	return dto.AuthProfileOut{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

// Refresh refresh user authentication token.
func (u *usecase) Refresh(ctx context.Context, payload *dto.AuthRefreshIn) (dto.AuthRefreshOut, error) {
	auth, err := u.authRepository.FindByToken(ctx, payload.RefreshToken)
	if err != nil {
		return dto.AuthRefreshOut{}, err
	}
	if err := auth.VerifyTokenExpires(); err != nil {
		return dto.AuthRefreshOut{}, err
	}

	accessToken, _, err := u.jwtProvider.GenerateAccessToken(auth.UserID)
	if err != nil {
		return dto.AuthRefreshOut{}, err
	}
	refreshToken, expires, err := u.jwtProvider.GenerateRefreshToken(auth.UserID)
	if err != nil {
		return dto.AuthRefreshOut{}, err
	}

	if err := u.authRepository.DeleteByToken(ctx, auth.Token); err != nil {
		return dto.AuthRefreshOut{}, err
	}

	newAuth := &entity.Auth{UserID: auth.UserID, Token: refreshToken, ExpiresAt: expires}
	if err := u.authRepository.Store(ctx, newAuth); err != nil {
		return dto.AuthRefreshOut{}, err
	}

	return dto.AuthRefreshOut{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
