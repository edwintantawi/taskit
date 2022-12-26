// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/edwintantawi/taskit/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// AuthUsecase is an autogenerated mock type for the AuthUsecase type
type AuthUsecase struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, payload
func (_m *AuthUsecase) Login(ctx context.Context, payload *domain.LoginAuthIn) (domain.LoginAuthOut, error) {
	ret := _m.Called(ctx, payload)

	var r0 domain.LoginAuthOut
	if rf, ok := ret.Get(0).(func(context.Context, *domain.LoginAuthIn) domain.LoginAuthOut); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(domain.LoginAuthOut)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.LoginAuthIn) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthUsecase creates a new instance of AuthUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthUsecase(t mockConstructorTestingTNewAuthUsecase) *AuthUsecase {
	mock := &AuthUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}