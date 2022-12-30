// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/edwintantawi/taskit/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// AuthRepository is an autogenerated mock type for the AuthRepository type
type AuthRepository struct {
	mock.Mock
}

// DeleteByToken provides a mock function with given fields: ctx, token
func (_m *AuthRepository) DeleteByToken(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByToken provides a mock function with given fields: ctx, token
func (_m *AuthRepository) FindByToken(ctx context.Context, token string) (entity.Auth, error) {
	ret := _m.Called(ctx, token)

	var r0 entity.Auth
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.Auth); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(entity.Auth)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, a
func (_m *AuthRepository) Store(ctx context.Context, a *entity.Auth) error {
	ret := _m.Called(ctx, a)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Auth) error); ok {
		r0 = rf(ctx, a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyAvailableByToken provides a mock function with given fields: ctx, token
func (_m *AuthRepository) VerifyAvailableByToken(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAuthRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthRepository creates a new instance of AuthRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthRepository(t mockConstructorTestingTNewAuthRepository) *AuthRepository {
	mock := &AuthRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
