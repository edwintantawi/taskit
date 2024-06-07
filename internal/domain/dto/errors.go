package dto

import "errors"

var (
	ErrEmailEmpty        = errors.New("dto.email_empty")
	ErrPasswordEmpty     = errors.New("dto.password_empty")
	ErrNameEmpty         = errors.New("dto.name_empty")
	ErrRefreshTokenEmpty = errors.New("dto.refresh_token_empty")
	ErrContentEmpty      = errors.New("dto.content_empty")
	ErrTitleEmpty        = errors.New("dto.title_empty")
)
