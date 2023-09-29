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

// CheckFriends provides a mock function with given fields: t
func (_m *RepositoryI) CheckFriends(t *models.FriendRelation) (bool, error) {
	ret := _m.Called(t)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.FriendRelation) (bool, error)); ok {
		return rf(t)
	}
	if rf, ok := ret.Get(0).(func(*models.FriendRelation) bool); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.FriendRelation) error); ok {
		r1 = rf(t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateFriendRelation provides a mock function with given fields: t
func (_m *RepositoryI) CreateFriendRelation(t *models.FriendRelation) error {
	ret := _m.Called(t)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.FriendRelation) error); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFriendRelation provides a mock function with given fields: friendRel
func (_m *RepositoryI) DeleteFriendRelation(friendRel *models.FriendRelation) error {
	ret := _m.Called(friendRel)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.FriendRelation) error); ok {
		r0 = rf(friendRel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserFriends provides a mock function with given fields: userID
func (_m *RepositoryI) GetUserFriends(userID uint64) ([]uint64, error) {
	ret := _m.Called(userID)

	var r0 []uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) ([]uint64, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uint64) []uint64); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserSubs provides a mock function with given fields: userID
func (_m *RepositoryI) GetUserSubs(userID uint64) ([]uint64, error) {
	ret := _m.Called(userID)

	var r0 []uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) ([]uint64, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uint64) []uint64); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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