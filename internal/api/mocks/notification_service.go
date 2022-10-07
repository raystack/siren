// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	notification "github.com/odpf/siren/core/notification"
	mock "github.com/stretchr/testify/mock"
)

// NotificationService is an autogenerated mock type for the NotificationService type
type NotificationService struct {
	mock.Mock
}

type NotificationService_Expecter struct {
	mock *mock.Mock
}

func (_m *NotificationService) EXPECT() *NotificationService_Expecter {
	return &NotificationService_Expecter{mock: &_m.Mock}
}

// DispatchBySubscription provides a mock function with given fields: ctx, n
func (_m *NotificationService) DispatchBySubscription(ctx context.Context, n notification.Notification) error {
	ret := _m.Called(ctx, n)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, notification.Notification) error); ok {
		r0 = rf(ctx, n)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NotificationService_DispatchBySubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DispatchBySubscription'
type NotificationService_DispatchBySubscription_Call struct {
	*mock.Call
}

// DispatchBySubscription is a helper method to define mock.On call
//  - ctx context.Context
//  - n notification.Notification
func (_e *NotificationService_Expecter) DispatchBySubscription(ctx interface{}, n interface{}) *NotificationService_DispatchBySubscription_Call {
	return &NotificationService_DispatchBySubscription_Call{Call: _e.mock.On("DispatchBySubscription", ctx, n)}
}

func (_c *NotificationService_DispatchBySubscription_Call) Run(run func(ctx context.Context, n notification.Notification)) *NotificationService_DispatchBySubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(notification.Notification))
	})
	return _c
}

func (_c *NotificationService_DispatchBySubscription_Call) Return(_a0 error) *NotificationService_DispatchBySubscription_Call {
	_c.Call.Return(_a0)
	return _c
}

// DispatchDirect provides a mock function with given fields: ctx, n, receiverID
func (_m *NotificationService) DispatchDirect(ctx context.Context, n notification.Notification, receiverID uint64) error {
	ret := _m.Called(ctx, n, receiverID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, notification.Notification, uint64) error); ok {
		r0 = rf(ctx, n, receiverID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NotificationService_DispatchDirect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DispatchDirect'
type NotificationService_DispatchDirect_Call struct {
	*mock.Call
}

// DispatchDirect is a helper method to define mock.On call
//  - ctx context.Context
//  - n notification.Notification
//  - receiverID uint64
func (_e *NotificationService_Expecter) DispatchDirect(ctx interface{}, n interface{}, receiverID interface{}) *NotificationService_DispatchDirect_Call {
	return &NotificationService_DispatchDirect_Call{Call: _e.mock.On("DispatchDirect", ctx, n, receiverID)}
}

func (_c *NotificationService_DispatchDirect_Call) Run(run func(ctx context.Context, n notification.Notification, receiverID uint64)) *NotificationService_DispatchDirect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(notification.Notification), args[2].(uint64))
	})
	return _c
}

func (_c *NotificationService_DispatchDirect_Call) Return(_a0 error) *NotificationService_DispatchDirect_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewNotificationService interface {
	mock.TestingT
	Cleanup(func())
}

// NewNotificationService creates a new instance of NotificationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNotificationService(t mockConstructorTestingTNewNotificationService) *NotificationService {
	mock := &NotificationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
