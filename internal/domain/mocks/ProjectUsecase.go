// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/edwintantawi/taskit/internal/domain/dto"

	mock "github.com/stretchr/testify/mock"
)

// ProjectUsecase is an autogenerated mock type for the ProjectUsecase type
type ProjectUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, payload
func (_m *ProjectUsecase) Create(ctx context.Context, payload *dto.ProjectCreateIn) (dto.ProjectCreateOut, error) {
	ret := _m.Called(ctx, payload)

	var r0 dto.ProjectCreateOut
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ProjectCreateIn) dto.ProjectCreateOut); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Get(0).(dto.ProjectCreateOut)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dto.ProjectCreateIn) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProjectUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewProjectUsecase creates a new instance of ProjectUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProjectUsecase(t mockConstructorTestingTNewProjectUsecase) *ProjectUsecase {
	mock := &ProjectUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}