// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"

	provider "github.com/goto/siren/core/provider"
	mock "github.com/stretchr/testify/mock"
)

// ProviderService is an autogenerated mock type for the ProviderService type
type ProviderService struct {
	mock.Mock
}

type ProviderService_Expecter struct {
	mock *mock.Mock
}

func (_m *ProviderService) EXPECT() *ProviderService_Expecter {
	return &ProviderService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *ProviderService) Create(_a0 context.Context, _a1 *provider.Provider) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *provider.Provider) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProviderService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type ProviderService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *provider.Provider
func (_e *ProviderService_Expecter) Create(_a0 interface{}, _a1 interface{}) *ProviderService_Create_Call {
	return &ProviderService_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *ProviderService_Create_Call) Run(run func(_a0 context.Context, _a1 *provider.Provider)) *ProviderService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*provider.Provider))
	})
	return _c
}

func (_c *ProviderService_Create_Call) Return(_a0 error) *ProviderService_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ProviderService_Create_Call) RunAndReturn(run func(context.Context, *provider.Provider) error) *ProviderService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *ProviderService) Delete(_a0 context.Context, _a1 uint64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProviderService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type ProviderService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 uint64
func (_e *ProviderService_Expecter) Delete(_a0 interface{}, _a1 interface{}) *ProviderService_Delete_Call {
	return &ProviderService_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *ProviderService_Delete_Call) Run(run func(_a0 context.Context, _a1 uint64)) *ProviderService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *ProviderService_Delete_Call) Return(_a0 error) *ProviderService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ProviderService_Delete_Call) RunAndReturn(run func(context.Context, uint64) error) *ProviderService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *ProviderService) Get(_a0 context.Context, _a1 uint64) (*provider.Provider, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *provider.Provider
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*provider.Provider, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *provider.Provider); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*provider.Provider)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProviderService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ProviderService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 uint64
func (_e *ProviderService_Expecter) Get(_a0 interface{}, _a1 interface{}) *ProviderService_Get_Call {
	return &ProviderService_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *ProviderService_Get_Call) Run(run func(_a0 context.Context, _a1 uint64)) *ProviderService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *ProviderService_Get_Call) Return(_a0 *provider.Provider, _a1 error) *ProviderService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderService_Get_Call) RunAndReturn(run func(context.Context, uint64) (*provider.Provider, error)) *ProviderService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: _a0, _a1
func (_m *ProviderService) List(_a0 context.Context, _a1 provider.Filter) ([]provider.Provider, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []provider.Provider
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, provider.Filter) ([]provider.Provider, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, provider.Filter) []provider.Provider); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]provider.Provider)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, provider.Filter) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProviderService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ProviderService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 provider.Filter
func (_e *ProviderService_Expecter) List(_a0 interface{}, _a1 interface{}) *ProviderService_List_Call {
	return &ProviderService_List_Call{Call: _e.mock.On("List", _a0, _a1)}
}

func (_c *ProviderService_List_Call) Run(run func(_a0 context.Context, _a1 provider.Filter)) *ProviderService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(provider.Filter))
	})
	return _c
}

func (_c *ProviderService_List_Call) Return(_a0 []provider.Provider, _a1 error) *ProviderService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderService_List_Call) RunAndReturn(run func(context.Context, provider.Filter) ([]provider.Provider, error)) *ProviderService_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *ProviderService) Update(_a0 context.Context, _a1 *provider.Provider) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *provider.Provider) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProviderService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ProviderService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *provider.Provider
func (_e *ProviderService_Expecter) Update(_a0 interface{}, _a1 interface{}) *ProviderService_Update_Call {
	return &ProviderService_Update_Call{Call: _e.mock.On("Update", _a0, _a1)}
}

func (_c *ProviderService_Update_Call) Run(run func(_a0 context.Context, _a1 *provider.Provider)) *ProviderService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*provider.Provider))
	})
	return _c
}

func (_c *ProviderService_Update_Call) Return(_a0 error) *ProviderService_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ProviderService_Update_Call) RunAndReturn(run func(context.Context, *provider.Provider) error) *ProviderService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewProviderService creates a new instance of ProviderService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProviderService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProviderService {
	mock := &ProviderService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
