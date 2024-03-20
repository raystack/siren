// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	notification "github.com/goto/siren/core/notification"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// Commit provides a mock function with given fields: ctx
func (_m *Repository) Commit(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type Repository_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Repository_Expecter) Commit(ctx interface{}) *Repository_Commit_Call {
	return &Repository_Commit_Call{Call: _e.mock.On("Commit", ctx)}
}

func (_c *Repository_Commit_Call) Run(run func(ctx context.Context)) *Repository_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Repository_Commit_Call) Return(_a0 error) *Repository_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Commit_Call) RunAndReturn(run func(context.Context) error) *Repository_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *Repository) Create(_a0 context.Context, _a1 notification.Notification) (notification.Notification, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 notification.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, notification.Notification) (notification.Notification, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, notification.Notification) notification.Notification); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(notification.Notification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, notification.Notification) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Repository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 notification.Notification
func (_e *Repository_Expecter) Create(_a0 interface{}, _a1 interface{}) *Repository_Create_Call {
	return &Repository_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *Repository_Create_Call) Run(run func(_a0 context.Context, _a1 notification.Notification)) *Repository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(notification.Notification))
	})
	return _c
}

func (_c *Repository_Create_Call) Return(_a0 notification.Notification, _a1 error) *Repository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Create_Call) RunAndReturn(run func(context.Context, notification.Notification) (notification.Notification, error)) *Repository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with given fields: ctx, err
func (_m *Repository) Rollback(ctx context.Context, err error) error {
	ret := _m.Called(ctx, err)

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, error) error); ok {
		r0 = rf(ctx, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type Repository_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *Repository_Expecter) Rollback(ctx interface{}, err interface{}) *Repository_Rollback_Call {
	return &Repository_Rollback_Call{Call: _e.mock.On("Rollback", ctx, err)}
}

func (_c *Repository_Rollback_Call) Run(run func(ctx context.Context, err error)) *Repository_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *Repository_Rollback_Call) Return(_a0 error) *Repository_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Rollback_Call) RunAndReturn(run func(context.Context, error) error) *Repository_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// WithTransaction provides a mock function with given fields: ctx
func (_m *Repository) WithTransaction(ctx context.Context) context.Context {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

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

// Repository_WithTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTransaction'
type Repository_WithTransaction_Call struct {
	*mock.Call
}

// WithTransaction is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Repository_Expecter) WithTransaction(ctx interface{}) *Repository_WithTransaction_Call {
	return &Repository_WithTransaction_Call{Call: _e.mock.On("WithTransaction", ctx)}
}

func (_c *Repository_WithTransaction_Call) Run(run func(ctx context.Context)) *Repository_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Repository_WithTransaction_Call) Return(_a0 context.Context) *Repository_WithTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_WithTransaction_Call) RunAndReturn(run func(context.Context) context.Context) *Repository_WithTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}