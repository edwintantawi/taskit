package usecase

import (
	"context"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type Usecase struct {
	validator      domain.ValidatorProvider
	authRepository domain.AuthRepository
	userRepository domain.UserRepository
	hashProvider   domain.HashProvider
	jwtProvider    domain.JWTProvider
}

// New create a new auth usecase.
func New(
	validator domain.ValidatorProvider,
	authRepository domain.AuthRepository,
	userRepository domain.UserRepository,
	hashProvider domain.HashProvider,
	jwtProvider domain.JWTProvider,
) Usecase {
	return Usecase{
		validator:      validator,
		authRepository: authRepository,
		userRepository: userRepository,
		hashProvider:   hashProvider,
		jwtProvider:    jwtProvider,
	}
}

// Login authenticates a user.
func (u *Usecase) Login(ctx context.Context, payload *dto.AuthLoginIn) (dto.AuthLoginOut, error) {
	user := entity.User{Email: payload.Email, Password: payload.Password}
	if err := u.validator.Validate(&user); err != nil {
		return dto.AuthLoginOut{}, err
	}

	targetUser, err := u.userRepository.FindByEmail(ctx, user.Email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return dto.AuthLoginOut{}, domain.ErrEmailNotExist
	} else if err != nil {
		return dto.AuthLoginOut{}, err
	}

	if err := u.hashProvider.Compare(user.Password, targetUser.Password); err != nil {
		return dto.AuthLoginOut{}, domain.ErrPasswordIncorrect
	}

	accessToken, _, err := u.jwtProvider.GenerateAccessToken(targetUser.ID)
	if err != nil {
		return dto.AuthLoginOut{}, err
	}
	refreshToken, expires, err := u.jwtProvider.GenerateRefreshToken(targetUser.ID)
	if err != nil {
		return dto.AuthLoginOut{}, err
	}

	auth := &entity.Auth{UserID: targetUser.ID, Token: refreshToken, ExpiresAt: expires}
	if err := u.authRepository.Store(ctx, auth); err != nil {
		return dto.AuthLoginOut{}, err
	}

	return dto.AuthLoginOut{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Logout remove user authentication.
func (u *Usecase) Logout(ctx context.Context, payload *dto.AuthLogoutIn) error {
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
func (u *Usecase) GetProfile(ctx context.Context, payload *dto.AuthProfileIn) (dto.AuthProfileOut, error) {
	user, err := u.userRepository.FindByID(ctx, payload.UserID)
	if err != nil {
		return dto.AuthProfileOut{}, err
	}
	return dto.AuthProfileOut{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

// Refresh refresh user authentication token.
func (u *Usecase) Refresh(ctx context.Context, payload *dto.AuthRefreshIn) (dto.AuthRefreshOut, error) {
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
