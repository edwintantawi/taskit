package domain

import "context"

type UserUsecase interface {
	Create(ctx context.Context, payload *CreateUserIn) (CreateUserOut, error)
}
