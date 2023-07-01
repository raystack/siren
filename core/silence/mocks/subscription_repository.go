// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	silence "github.com/raystack/siren/core/silence"
	mock "github.com/stretchr/testify/mock"
)

// SubscriptionRepository is an autogenerated mock type for the Repository type
type SubscriptionRepository struct {
	mock.Mock
}

type SubscriptionRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *SubscriptionRepository) EXPECT() *SubscriptionRepository_Expecter {
	return &SubscriptionRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *SubscriptionRepository) Create(_a0 context.Context, _a1 silence.Silence) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, silence.Silence) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, silence.Silence) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type SubscriptionRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 silence.Silence
func (_e *SubscriptionRepository_Expecter) Create(_a0 interface{}, _a1 interface{}) *SubscriptionRepository_Create_Call {
	return &SubscriptionRepository_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *SubscriptionRepository_Create_Call) Run(run func(_a0 context.Context, _a1 silence.Silence)) *SubscriptionRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(silence.Silence))
	})
	return _c
}

func (_c *SubscriptionRepository_Create_Call) Return(_a0 string, _a1 error) *SubscriptionRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *SubscriptionRepository) Get(ctx context.Context, id string) (silence.Silence, error) {
	ret := _m.Called(ctx, id)

	var r0 silence.Silence
	if rf, ok := ret.Get(0).(func(context.Context, string) silence.Silence); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(silence.Silence)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type SubscriptionRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *SubscriptionRepository_Expecter) Get(ctx interface{}, id interface{}) *SubscriptionRepository_Get_Call {
	return &SubscriptionRepository_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *SubscriptionRepository_Get_Call) Run(run func(ctx context.Context, id string)) *SubscriptionRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SubscriptionRepository_Get_Call) Return(_a0 silence.Silence, _a1 error) *SubscriptionRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// List provides a mock function with given fields: _a0, _a1
func (_m *SubscriptionRepository) List(_a0 context.Context, _a1 silence.Filter) ([]silence.Silence, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []silence.Silence
	if rf, ok := ret.Get(0).(func(context.Context, silence.Filter) []silence.Silence); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]silence.Silence)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, silence.Filter) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type SubscriptionRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 silence.Filter
func (_e *SubscriptionRepository_Expecter) List(_a0 interface{}, _a1 interface{}) *SubscriptionRepository_List_Call {
	return &SubscriptionRepository_List_Call{Call: _e.mock.On("List", _a0, _a1)}
}

func (_c *SubscriptionRepository_List_Call) Run(run func(_a0 context.Context, _a1 silence.Filter)) *SubscriptionRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(silence.Filter))
	})
	return _c
}

func (_c *SubscriptionRepository_List_Call) Return(_a0 []silence.Silence, _a1 error) *SubscriptionRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// SoftDelete provides a mock function with given fields: ctx, id
func (_m *SubscriptionRepository) SoftDelete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscriptionRepository_SoftDelete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SoftDelete'
type SubscriptionRepository_SoftDelete_Call struct {
	*mock.Call
}

// SoftDelete is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *SubscriptionRepository_Expecter) SoftDelete(ctx interface{}, id interface{}) *SubscriptionRepository_SoftDelete_Call {
	return &SubscriptionRepository_SoftDelete_Call{Call: _e.mock.On("SoftDelete", ctx, id)}
}

func (_c *SubscriptionRepository_SoftDelete_Call) Run(run func(ctx context.Context, id string)) *SubscriptionRepository_SoftDelete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SubscriptionRepository_SoftDelete_Call) Return(_a0 error) *SubscriptionRepository_SoftDelete_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewSubscriptionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSubscriptionRepository(t mockConstructorTestingTNewSubscriptionRepository) *SubscriptionRepository {
	mock := &SubscriptionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
