// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/edwintantawi/taskit/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// TaskUsecase is an autogenerated mock type for the TaskUsecase type
type TaskUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, payload
func (_m *TaskUsecase) Create(ctx context.Context, payload *domain.CreateTaskIn) (domain.CreateTaskOut, error) {
	ret := _m.Called(ctx, payload)

	var r0 domain.CreateTaskOut
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CreateTaskIn) domain.CreateTaskOut); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(domain.CreateTaskOut)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.CreateTaskIn) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTaskUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewTaskUsecase creates a new instance of TaskUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTaskUsecase(t mockConstructorTestingTNewTaskUsecase) *TaskUsecase {
	mock := &TaskUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}