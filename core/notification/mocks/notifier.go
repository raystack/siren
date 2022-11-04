// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	notification "github.com/odpf/siren/core/notification"
	mock "github.com/stretchr/testify/mock"
)

// Notifier is an autogenerated mock type for the Notifier type
type Notifier struct {
	mock.Mock
}

type Notifier_Expecter struct {
	mock *mock.Mock
}

func (_m *Notifier) EXPECT() *Notifier_Expecter {
	return &Notifier_Expecter{mock: &_m.Mock}
}

// GetSystemDefaultTemplate provides a mock function with given fields:
func (_m *Notifier) GetSystemDefaultTemplate() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Notifier_GetSystemDefaultTemplate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSystemDefaultTemplate'
type Notifier_GetSystemDefaultTemplate_Call struct {
	*mock.Call
}

// GetSystemDefaultTemplate is a helper method to define mock.On call
func (_e *Notifier_Expecter) GetSystemDefaultTemplate() *Notifier_GetSystemDefaultTemplate_Call {
	return &Notifier_GetSystemDefaultTemplate_Call{Call: _e.mock.On("GetSystemDefaultTemplate")}
}

func (_c *Notifier_GetSystemDefaultTemplate_Call) Run(run func()) *Notifier_GetSystemDefaultTemplate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Notifier_GetSystemDefaultTemplate_Call) Return(_a0 string) *Notifier_GetSystemDefaultTemplate_Call {
	_c.Call.Return(_a0)
	return _c
}

// PostHookQueueTransformConfigs provides a mock function with given fields: ctx, notificationConfigMap
func (_m *Notifier) PostHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(ctx, notificationConfigMap)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(ctx, notificationConfigMap)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, notificationConfigMap)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Notifier_PostHookQueueTransformConfigs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PostHookQueueTransformConfigs'
type Notifier_PostHookQueueTransformConfigs_Call struct {
	*mock.Call
}

// PostHookQueueTransformConfigs is a helper method to define mock.On call
//  - ctx context.Context
//  - notificationConfigMap map[string]interface{}
func (_e *Notifier_Expecter) PostHookQueueTransformConfigs(ctx interface{}, notificationConfigMap interface{}) *Notifier_PostHookQueueTransformConfigs_Call {
	return &Notifier_PostHookQueueTransformConfigs_Call{Call: _e.mock.On("PostHookQueueTransformConfigs", ctx, notificationConfigMap)}
}

func (_c *Notifier_PostHookQueueTransformConfigs_Call) Run(run func(ctx context.Context, notificationConfigMap map[string]interface{})) *Notifier_PostHookQueueTransformConfigs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *Notifier_PostHookQueueTransformConfigs_Call) Return(_a0 map[string]interface{}, _a1 error) *Notifier_PostHookQueueTransformConfigs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// PreHookQueueTransformConfigs provides a mock function with given fields: ctx, notificationConfigMap
func (_m *Notifier) PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(ctx, notificationConfigMap)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(ctx, notificationConfigMap)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, notificationConfigMap)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Notifier_PreHookQueueTransformConfigs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PreHookQueueTransformConfigs'
type Notifier_PreHookQueueTransformConfigs_Call struct {
	*mock.Call
}

// PreHookQueueTransformConfigs is a helper method to define mock.On call
//  - ctx context.Context
//  - notificationConfigMap map[string]interface{}
func (_e *Notifier_Expecter) PreHookQueueTransformConfigs(ctx interface{}, notificationConfigMap interface{}) *Notifier_PreHookQueueTransformConfigs_Call {
	return &Notifier_PreHookQueueTransformConfigs_Call{Call: _e.mock.On("PreHookQueueTransformConfigs", ctx, notificationConfigMap)}
}

func (_c *Notifier_PreHookQueueTransformConfigs_Call) Run(run func(ctx context.Context, notificationConfigMap map[string]interface{})) *Notifier_PreHookQueueTransformConfigs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]interface{}))
	})
	return _c
}

func (_c *Notifier_PreHookQueueTransformConfigs_Call) Return(_a0 map[string]interface{}, _a1 error) *Notifier_PreHookQueueTransformConfigs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Send provides a mock function with given fields: ctx, message
func (_m *Notifier) Send(ctx context.Context, message notification.Message) (bool, error) {
	ret := _m.Called(ctx, message)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, notification.Message) bool); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, notification.Message) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Notifier_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type Notifier_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//  - ctx context.Context
//  - message notification.Message
func (_e *Notifier_Expecter) Send(ctx interface{}, message interface{}) *Notifier_Send_Call {
	return &Notifier_Send_Call{Call: _e.mock.On("Send", ctx, message)}
}

func (_c *Notifier_Send_Call) Run(run func(ctx context.Context, message notification.Message)) *Notifier_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(notification.Message))
	})
	return _c
}

func (_c *Notifier_Send_Call) Return(_a0 bool, _a1 error) *Notifier_Send_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewNotifier interface {
	mock.TestingT
	Cleanup(func())
}

// NewNotifier creates a new instance of Notifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNotifier(t mockConstructorTestingTNewNotifier) *Notifier {
	mock := &Notifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}