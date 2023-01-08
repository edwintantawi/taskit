// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/edwintantawi/taskit/internal/domain/dto"

	mock "github.com/stretchr/testify/mock"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, payload
func (_m *UserUsecase) Create(ctx context.Context, payload *dto.UserCreateIn) (dto.UserCreateOut, error) {
	ret := _m.Called(ctx, payload)

	var r0 dto.UserCreateOut
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UserCreateIn) dto.UserCreateOut); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(dto.UserCreateOut)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dto.UserCreateIn) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserUsecase creates a new instance of UserUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserUsecase(t mockConstructorTestingTNewUserUsecase) *UserUsecase {
	mock := &UserUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
