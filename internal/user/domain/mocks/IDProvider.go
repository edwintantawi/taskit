// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IDProvider is an autogenerated mock type for the IDProvider type
type IDProvider struct {
	mock.Mock
}

// Generate provides a mock function with given fields:
func (_m *IDProvider) Generate() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewIDProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewIDProvider creates a new instance of IDProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIDProvider(t mockConstructorTestingTNewIDProvider) *IDProvider {
	mock := &IDProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
