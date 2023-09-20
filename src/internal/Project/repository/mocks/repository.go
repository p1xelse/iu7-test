// Code generated by mockery v2.23.2. DO NOT EDIT.

package mocks

import (
	models "timetracker/models"

	mock "github.com/stretchr/testify/mock"
)

// RepositoryI is an autogenerated mock type for the RepositoryI type
type RepositoryI struct {
	mock.Mock
}

// CreateProject provides a mock function with given fields: e
func (_m *RepositoryI) CreateProject(e *models.Project) error {
	ret := _m.Called(e)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Project) error); ok {
		r0 = rf(e)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteProject provides a mock function with given fields: id
func (_m *RepositoryI) DeleteProject(id uint64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProject provides a mock function with given fields: id
func (_m *RepositoryI) GetProject(id uint64) (*models.Project, error) {
	ret := _m.Called(id)

	var r0 *models.Project
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*models.Project, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) *models.Project); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserProjects provides a mock function with given fields: userID
func (_m *RepositoryI) GetUserProjects(userID uint64) ([]*models.Project, error) {
	ret := _m.Called(userID)

	var r0 []*models.Project
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) ([]*models.Project, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uint64) []*models.Project); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProject provides a mock function with given fields: e
func (_m *RepositoryI) UpdateProject(e *models.Project) error {
	ret := _m.Called(e)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Project) error); ok {
		r0 = rf(e)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepositoryI interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepositoryI creates a new instance of RepositoryI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepositoryI(t mockConstructorTestingTNewRepositoryI) *RepositoryI {
	mock := &RepositoryI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
