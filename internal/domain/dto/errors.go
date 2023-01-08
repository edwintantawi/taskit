package dto

import "errors"

var (
	ErrEmailEmpty    = errors.New("dto.email_empty")
	ErrPasswordEmpty = errors.New("dto.password_empty")
	ErrNameEmpty     = errors.New("dto.name_empty")

	ErrRefreshTokenEmpty = errors.New("dto.refresh_token_empty")

	ErrTaskContentEmpty = errors.New("dto.task_content_empty")
)
