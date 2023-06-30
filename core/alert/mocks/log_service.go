// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// LogService is an autogenerated mock type for the LogService type
type LogService struct {
	mock.Mock
}

type LogService_Expecter struct {
	mock *mock.Mock
}

func (_m *LogService) EXPECT() *LogService_Expecter {
	return &LogService_Expecter{mock: &_m.Mock}
}

// ListNotificationAlertIDsBySilenceID provides a mock function with given fields: ctx, silenceID
func (_m *LogService) ListNotificationAlertIDsBySilenceID(ctx context.Context, silenceID string) ([]int64, error) {
	ret := _m.Called(ctx, silenceID)

	var r0 []int64
	if rf, ok := ret.Get(0).(func(context.Context, string) []int64); ok {
		r0 = rf(ctx, silenceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int64)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, silenceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogService_ListNotificationAlertIDsBySilenceID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListNotificationAlertIDsBySilenceID'
type LogService_ListNotificationAlertIDsBySilenceID_Call struct {
	*mock.Call
}

// ListNotificationAlertIDsBySilenceID is a helper method to define mock.On call
//   - ctx context.Context
//   - silenceID string
func (_e *LogService_Expecter) ListNotificationAlertIDsBySilenceID(ctx interface{}, silenceID interface{}) *LogService_ListNotificationAlertIDsBySilenceID_Call {
	return &LogService_ListNotificationAlertIDsBySilenceID_Call{Call: _e.mock.On("ListNotificationAlertIDsBySilenceID", ctx, silenceID)}
}

func (_c *LogService_ListNotificationAlertIDsBySilenceID_Call) Run(run func(ctx context.Context, silenceID string)) *LogService_ListNotificationAlertIDsBySilenceID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *LogService_ListNotificationAlertIDsBySilenceID_Call) Return(_a0 []int64, _a1 error) *LogService_ListNotificationAlertIDsBySilenceID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewLogService interface {
	mock.TestingT
	Cleanup(func())
}

// NewLogService creates a new instance of LogService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLogService(t mockConstructorTestingTNewLogService) *LogService {
	mock := &LogService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}