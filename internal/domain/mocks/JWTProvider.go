// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/edwintantawi/taskit/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// JWTProvider is an autogenerated mock type for the JWTProvider type
type JWTProvider struct {
	mock.Mock
}

// GenerateAccessToken provides a mock function with given fields: userID
func (_m *JWTProvider) GenerateAccessToken(userID entity.UserID) (string, time.Time, error) {
	ret := _m.Called(userID)

	var r0 string
	if rf, ok := ret.Get(0).(func(entity.UserID) string); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 time.Time
	if rf, ok := ret.Get(1).(func(entity.UserID) time.Time); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(entity.UserID) error); ok {
		r2 = rf(userID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GenerateRefreshToken provides a mock function with given fields: userID
func (_m *JWTProvider) GenerateRefreshToken(userID entity.UserID) (string, time.Time, error) {
	ret := _m.Called(userID)

	var r0 string
	if rf, ok := ret.Get(0).(func(entity.UserID) string); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 time.Time
	if rf, ok := ret.Get(1).(func(entity.UserID) time.Time); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(entity.UserID) error); ok {
		r2 = rf(userID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// VerifyAccessToken provides a mock function with given fields: rawToken
func (_m *JWTProvider) VerifyAccessToken(rawToken string) (entity.UserID, error) {
	ret := _m.Called(rawToken)

	var r0 entity.UserID
	if rf, ok := ret.Get(0).(func(string) entity.UserID); ok {
		r0 = rf(rawToken)
	} else {
		r0 = ret.Get(0).(entity.UserID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(rawToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewJWTProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewJWTProvider creates a new instance of JWTProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJWTProvider(t mockConstructorTestingTNewJWTProvider) *JWTProvider {
	mock := &JWTProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
