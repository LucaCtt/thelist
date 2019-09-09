// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import common "github.com/LucaCtt/thelist/common"
import mock "github.com/stretchr/testify/mock"

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *Store) All() ([]*common.Item, error) {
	ret := _m.Called()

	var r0 []*common.Item
	if rf, ok := ret.Get(0).(func() []*common.Item); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*common.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *Store) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: item
func (_m *Store) Create(item *common.Item) error {
	ret := _m.Called(item)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.Item) error); ok {
		r0 = rf(item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *Store) Delete(id uint) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *Store) Get(id uint) (*common.Item, error) {
	ret := _m.Called(id)

	var r0 *common.Item
	if rf, ok := ret.Get(0).(func(uint) *common.Item); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetWatched provides a mock function with given fields: id, watched
func (_m *Store) SetWatched(id uint, watched bool) error {
	ret := _m.Called(id, watched)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, bool) error); ok {
		r0 = rf(id, watched)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}