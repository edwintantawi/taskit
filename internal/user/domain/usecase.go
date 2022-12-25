package domain

import "context"

// UserUsecase represent user usecase contract
type UserUsecase interface {
	Create(ctx context.Context, payload *CreateUserIn) (CreateUserOut, error)
}
