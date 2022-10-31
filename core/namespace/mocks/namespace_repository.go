// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	namespace "github.com/odpf/siren/core/namespace"
	mock "github.com/stretchr/testify/mock"
)

// NamespaceRepository is an autogenerated mock type for the Repository type
type NamespaceRepository struct {
	mock.Mock
}

type NamespaceRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *NamespaceRepository) EXPECT() *NamespaceRepository_Expecter {
	return &NamespaceRepository_Expecter{mock: &_m.Mock}
}

// Commit provides a mock function with given fields: ctx
func (_m *NamespaceRepository) Commit(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceRepository_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type NamespaceRepository_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
//  - ctx context.Context
func (_e *NamespaceRepository_Expecter) Commit(ctx interface{}) *NamespaceRepository_Commit_Call {
	return &NamespaceRepository_Commit_Call{Call: _e.mock.On("Commit", ctx)}
}

func (_c *NamespaceRepository_Commit_Call) Run(run func(ctx context.Context)) *NamespaceRepository_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *NamespaceRepository_Commit_Call) Return(_a0 error) *NamespaceRepository_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) Create(_a0 context.Context, _a1 *namespace.EncryptedNamespace) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *namespace.EncryptedNamespace) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type NamespaceRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 *namespace.EncryptedNamespace
func (_e *NamespaceRepository_Expecter) Create(_a0 interface{}, _a1 interface{}) *NamespaceRepository_Create_Call {
	return &NamespaceRepository_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *NamespaceRepository_Create_Call) Run(run func(_a0 context.Context, _a1 *namespace.EncryptedNamespace)) *NamespaceRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*namespace.EncryptedNamespace))
	})
	return _c
}

func (_c *NamespaceRepository_Create_Call) Return(_a0 error) *NamespaceRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) Delete(_a0 context.Context, _a1 uint64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type NamespaceRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 uint64
func (_e *NamespaceRepository_Expecter) Delete(_a0 interface{}, _a1 interface{}) *NamespaceRepository_Delete_Call {
	return &NamespaceRepository_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *NamespaceRepository_Delete_Call) Run(run func(_a0 context.Context, _a1 uint64)) *NamespaceRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *NamespaceRepository_Delete_Call) Return(_a0 error) *NamespaceRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) Get(_a0 context.Context, _a1 uint64) (*namespace.EncryptedNamespace, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *namespace.EncryptedNamespace
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *namespace.EncryptedNamespace); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*namespace.EncryptedNamespace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type NamespaceRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 uint64
func (_e *NamespaceRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *NamespaceRepository_Get_Call {
	return &NamespaceRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *NamespaceRepository_Get_Call) Run(run func(_a0 context.Context, _a1 uint64)) *NamespaceRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *NamespaceRepository_Get_Call) Return(_a0 *namespace.EncryptedNamespace, _a1 error) *NamespaceRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// List provides a mock function with given fields: _a0
func (_m *NamespaceRepository) List(_a0 context.Context) ([]namespace.EncryptedNamespace, error) {
	ret := _m.Called(_a0)

	var r0 []namespace.EncryptedNamespace
	if rf, ok := ret.Get(0).(func(context.Context) []namespace.EncryptedNamespace); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]namespace.EncryptedNamespace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespaceRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type NamespaceRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//  - _a0 context.Context
func (_e *NamespaceRepository_Expecter) List(_a0 interface{}) *NamespaceRepository_List_Call {
	return &NamespaceRepository_List_Call{Call: _e.mock.On("List", _a0)}
}

func (_c *NamespaceRepository_List_Call) Run(run func(_a0 context.Context)) *NamespaceRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *NamespaceRepository_List_Call) Return(_a0 []namespace.EncryptedNamespace, _a1 error) *NamespaceRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Rollback provides a mock function with given fields: ctx, err
func (_m *NamespaceRepository) Rollback(ctx context.Context, err error) error {
	ret := _m.Called(ctx, err)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, error) error); ok {
		r0 = rf(ctx, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceRepository_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type NamespaceRepository_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
//  - ctx context.Context
//  - err error
func (_e *NamespaceRepository_Expecter) Rollback(ctx interface{}, err interface{}) *NamespaceRepository_Rollback_Call {
	return &NamespaceRepository_Rollback_Call{Call: _e.mock.On("Rollback", ctx, err)}
}

func (_c *NamespaceRepository_Rollback_Call) Run(run func(ctx context.Context, err error)) *NamespaceRepository_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *NamespaceRepository_Rollback_Call) Return(_a0 error) *NamespaceRepository_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) Update(_a0 context.Context, _a1 *namespace.EncryptedNamespace) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *namespace.EncryptedNamespace) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type NamespaceRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 *namespace.EncryptedNamespace
func (_e *NamespaceRepository_Expecter) Update(_a0 interface{}, _a1 interface{}) *NamespaceRepository_Update_Call {
	return &NamespaceRepository_Update_Call{Call: _e.mock.On("Update", _a0, _a1)}
}

func (_c *NamespaceRepository_Update_Call) Run(run func(_a0 context.Context, _a1 *namespace.EncryptedNamespace)) *NamespaceRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*namespace.EncryptedNamespace))
	})
	return _c
}

func (_c *NamespaceRepository_Update_Call) Return(_a0 error) *NamespaceRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

// WithTransaction provides a mock function with given fields: ctx
func (_m *NamespaceRepository) WithTransaction(ctx context.Context) context.Context {
	ret := _m.Called(ctx)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(context.Context) context.Context); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// NamespaceRepository_WithTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTransaction'
type NamespaceRepository_WithTransaction_Call struct {
	*mock.Call
}

// WithTransaction is a helper method to define mock.On call
//  - ctx context.Context
func (_e *NamespaceRepository_Expecter) WithTransaction(ctx interface{}) *NamespaceRepository_WithTransaction_Call {
	return &NamespaceRepository_WithTransaction_Call{Call: _e.mock.On("WithTransaction", ctx)}
}

func (_c *NamespaceRepository_WithTransaction_Call) Run(run func(ctx context.Context)) *NamespaceRepository_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *NamespaceRepository_WithTransaction_Call) Return(_a0 context.Context) *NamespaceRepository_WithTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewNamespaceRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewNamespaceRepository creates a new instance of NamespaceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNamespaceRepository(t mockConstructorTestingTNewNamespaceRepository) *NamespaceRepository {
	mock := &NamespaceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
