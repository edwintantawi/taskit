package dto

import "errors"

var (
	ErrEmailEmpty    = errors.New("dto.email_empty")
	ErrPasswordEmpty = errors.New("dto.password_empty")
	ErrNameEmpty     = errors.New("dto.name_empty")
)
