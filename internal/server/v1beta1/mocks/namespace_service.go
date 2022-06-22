// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	namespace "github.com/odpf/siren/core/namespace"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// NamespaceService is an autogenerated mock type for the NamespaceService type
type NamespaceService struct {
	mock.Mock
}

type NamespaceService_Expecter struct {
	mock *mock.Mock
}

func (_m *NamespaceService) EXPECT() *NamespaceService_Expecter {
	return &NamespaceService_Expecter{mock: &_m.Mock}
}

// CreateNamespace provides a mock function with given fields: _a0
func (_m *NamespaceService) CreateNamespace(_a0 *namespace.Namespace) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*namespace.Namespace) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceService_CreateNamespace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateNamespace'
type NamespaceService_CreateNamespace_Call struct {
	*mock.Call
}

// CreateNamespace is a helper method to define mock.On call
//  - _a0 *namespace.Namespace
func (_e *NamespaceService_Expecter) CreateNamespace(_a0 interface{}) *NamespaceService_CreateNamespace_Call {
	return &NamespaceService_CreateNamespace_Call{Call: _e.mock.On("CreateNamespace", _a0)}
}

func (_c *NamespaceService_CreateNamespace_Call) Run(run func(_a0 *namespace.Namespace)) *NamespaceService_CreateNamespace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*namespace.Namespace))
	})
	return _c
}

func (_c *NamespaceService_CreateNamespace_Call) Return(_a0 error) *NamespaceService_CreateNamespace_Call {
	_c.Call.Return(_a0)
	return _c
}

// DeleteNamespace provides a mock function with given fields: _a0
func (_m *NamespaceService) DeleteNamespace(_a0 uint64) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceService_DeleteNamespace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteNamespace'
type NamespaceService_DeleteNamespace_Call struct {
	*mock.Call
}

// DeleteNamespace is a helper method to define mock.On call
//  - _a0 uint64
func (_e *NamespaceService_Expecter) DeleteNamespace(_a0 interface{}) *NamespaceService_DeleteNamespace_Call {
	return &NamespaceService_DeleteNamespace_Call{Call: _e.mock.On("DeleteNamespace", _a0)}
}

func (_c *NamespaceService_DeleteNamespace_Call) Run(run func(_a0 uint64)) *NamespaceService_DeleteNamespace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint64))
	})
	return _c
}

func (_c *NamespaceService_DeleteNamespace_Call) Return(_a0 error) *NamespaceService_DeleteNamespace_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetNamespace provides a mock function with given fields: _a0
func (_m *NamespaceService) GetNamespace(_a0 uint64) (*namespace.Namespace, error) {
	ret := _m.Called(_a0)

	var r0 *namespace.Namespace
	if rf, ok := ret.Get(0).(func(uint64) *namespace.Namespace); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*namespace.Namespace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceService_GetNamespace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNamespace'
type NamespaceService_GetNamespace_Call struct {
	*mock.Call
}

// GetNamespace is a helper method to define mock.On call
//  - _a0 uint64
func (_e *NamespaceService_Expecter) GetNamespace(_a0 interface{}) *NamespaceService_GetNamespace_Call {
	return &NamespaceService_GetNamespace_Call{Call: _e.mock.On("GetNamespace", _a0)}
}

func (_c *NamespaceService_GetNamespace_Call) Run(run func(_a0 uint64)) *NamespaceService_GetNamespace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint64))
	})
	return _c
}

func (_c *NamespaceService_GetNamespace_Call) Return(_a0 *namespace.Namespace, _a1 error) *NamespaceService_GetNamespace_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// ListNamespaces provides a mock function with given fields:
func (_m *NamespaceService) ListNamespaces() ([]*namespace.Namespace, error) {
	ret := _m.Called()

	var r0 []*namespace.Namespace
	if rf, ok := ret.Get(0).(func() []*namespace.Namespace); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*namespace.Namespace)
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

// NamespaceService_ListNamespaces_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListNamespaces'
type NamespaceService_ListNamespaces_Call struct {
	*mock.Call
}

// ListNamespaces is a helper method to define mock.On call
func (_e *NamespaceService_Expecter) ListNamespaces() *NamespaceService_ListNamespaces_Call {
	return &NamespaceService_ListNamespaces_Call{Call: _e.mock.On("ListNamespaces")}
}

func (_c *NamespaceService_ListNamespaces_Call) Run(run func()) *NamespaceService_ListNamespaces_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *NamespaceService_ListNamespaces_Call) Return(_a0 []*namespace.Namespace, _a1 error) *NamespaceService_ListNamespaces_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// UpdateNamespace provides a mock function with given fields: _a0
func (_m *NamespaceService) UpdateNamespace(_a0 *namespace.Namespace) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*namespace.Namespace) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceService_UpdateNamespace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateNamespace'
type NamespaceService_UpdateNamespace_Call struct {
	*mock.Call
}

// UpdateNamespace is a helper method to define mock.On call
//  - _a0 *namespace.Namespace
func (_e *NamespaceService_Expecter) UpdateNamespace(_a0 interface{}) *NamespaceService_UpdateNamespace_Call {
	return &NamespaceService_UpdateNamespace_Call{Call: _e.mock.On("UpdateNamespace", _a0)}
}

func (_c *NamespaceService_UpdateNamespace_Call) Run(run func(_a0 *namespace.Namespace)) *NamespaceService_UpdateNamespace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*namespace.Namespace))
	})
	return _c
}

func (_c *NamespaceService_UpdateNamespace_Call) Return(_a0 error) *NamespaceService_UpdateNamespace_Call {
	_c.Call.Return(_a0)
	return _c
}

// NewNamespaceService creates a new instance of NamespaceService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewNamespaceService(t testing.TB) *NamespaceService {
	mock := &NamespaceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}