// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	receiver "github.com/odpf/siren/core/receiver"
	mock "github.com/stretchr/testify/mock"
)

// ReceiverService is an autogenerated mock type for the ReceiverService type
type ReceiverService struct {
	mock.Mock
}

type ReceiverService_Expecter struct {
	mock *mock.Mock
}

func (_m *ReceiverService) EXPECT() *ReceiverService_Expecter {
	return &ReceiverService_Expecter{mock: &_m.Mock}
}

// BuildNotificationConfig provides a mock function with given fields: subsConfig, rcv
func (_m *ReceiverService) BuildNotificationConfig(subsConfig map[string]interface{}, rcv *receiver.Receiver) (map[string]interface{}, error) {
	ret := _m.Called(subsConfig, rcv)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(map[string]interface{}, *receiver.Receiver) map[string]interface{}); ok {
		r0 = rf(subsConfig, rcv)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(map[string]interface{}, *receiver.Receiver) error); ok {
		r1 = rf(subsConfig, rcv)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReceiverService_BuildNotificationConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BuildNotificationConfig'
type ReceiverService_BuildNotificationConfig_Call struct {
	*mock.Call
}

// BuildNotificationConfig is a helper method to define mock.On call
//  - subsConfig map[string]interface{}
//  - rcv *receiver.Receiver
func (_e *ReceiverService_Expecter) BuildNotificationConfig(subsConfig interface{}, rcv interface{}) *ReceiverService_BuildNotificationConfig_Call {
	return &ReceiverService_BuildNotificationConfig_Call{Call: _e.mock.On("BuildNotificationConfig", subsConfig, rcv)}
}

func (_c *ReceiverService_BuildNotificationConfig_Call) Run(run func(subsConfig map[string]interface{}, rcv *receiver.Receiver)) *ReceiverService_BuildNotificationConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]interface{}), args[1].(*receiver.Receiver))
	})
	return _c
}

func (_c *ReceiverService_BuildNotificationConfig_Call) Return(_a0 map[string]interface{}, _a1 error) *ReceiverService_BuildNotificationConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Create provides a mock function with given fields: ctx, rcv
func (_m *ReceiverService) Create(ctx context.Context, rcv *receiver.Receiver) error {
	ret := _m.Called(ctx, rcv)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *receiver.Receiver) error); ok {
		r0 = rf(ctx, rcv)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReceiverService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type ReceiverService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//  - ctx context.Context
//  - rcv *receiver.Receiver
func (_e *ReceiverService_Expecter) Create(ctx interface{}, rcv interface{}) *ReceiverService_Create_Call {
	return &ReceiverService_Create_Call{Call: _e.mock.On("Create", ctx, rcv)}
}

func (_c *ReceiverService_Create_Call) Run(run func(ctx context.Context, rcv *receiver.Receiver)) *ReceiverService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*receiver.Receiver))
	})
	return _c
}

func (_c *ReceiverService_Create_Call) Return(_a0 error) *ReceiverService_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *ReceiverService) Delete(ctx context.Context, id uint64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReceiverService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type ReceiverService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//  - ctx context.Context
//  - id uint64
func (_e *ReceiverService_Expecter) Delete(ctx interface{}, id interface{}) *ReceiverService_Delete_Call {
	return &ReceiverService_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *ReceiverService_Delete_Call) Run(run func(ctx context.Context, id uint64)) *ReceiverService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *ReceiverService_Delete_Call) Return(_a0 error) *ReceiverService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *ReceiverService) Get(ctx context.Context, id uint64) (*receiver.Receiver, error) {
	ret := _m.Called(ctx, id)

	var r0 *receiver.Receiver
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *receiver.Receiver); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*receiver.Receiver)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReceiverService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ReceiverService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//  - ctx context.Context
//  - id uint64
func (_e *ReceiverService_Expecter) Get(ctx interface{}, id interface{}) *ReceiverService_Get_Call {
	return &ReceiverService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *ReceiverService_Get_Call) Run(run func(ctx context.Context, id uint64)) *ReceiverService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *ReceiverService_Get_Call) Return(_a0 *receiver.Receiver, _a1 error) *ReceiverService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// List provides a mock function with given fields: ctx, flt
func (_m *ReceiverService) List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error) {
	ret := _m.Called(ctx, flt)

	var r0 []receiver.Receiver
	if rf, ok := ret.Get(0).(func(context.Context, receiver.Filter) []receiver.Receiver); ok {
		r0 = rf(ctx, flt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]receiver.Receiver)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, receiver.Filter) error); ok {
		r1 = rf(ctx, flt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReceiverService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ReceiverService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//  - ctx context.Context
//  - flt receiver.Filter
func (_e *ReceiverService_Expecter) List(ctx interface{}, flt interface{}) *ReceiverService_List_Call {
	return &ReceiverService_List_Call{Call: _e.mock.On("List", ctx, flt)}
}

func (_c *ReceiverService_List_Call) Run(run func(ctx context.Context, flt receiver.Filter)) *ReceiverService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(receiver.Filter))
	})
	return _c
}

func (_c *ReceiverService_List_Call) Return(_a0 []receiver.Receiver, _a1 error) *ReceiverService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Notify provides a mock function with given fields: ctx, id, payloadMessage
func (_m *ReceiverService) Notify(ctx context.Context, id uint64, payloadMessage map[string]interface{}) error {
	ret := _m.Called(ctx, id, payloadMessage)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, map[string]interface{}) error); ok {
		r0 = rf(ctx, id, payloadMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReceiverService_Notify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Notify'
type ReceiverService_Notify_Call struct {
	*mock.Call
}

// Notify is a helper method to define mock.On call
//  - ctx context.Context
//  - id uint64
//  - payloadMessage map[string]interface{}
func (_e *ReceiverService_Expecter) Notify(ctx interface{}, id interface{}, payloadMessage interface{}) *ReceiverService_Notify_Call {
	return &ReceiverService_Notify_Call{Call: _e.mock.On("Notify", ctx, id, payloadMessage)}
}

func (_c *ReceiverService_Notify_Call) Run(run func(ctx context.Context, id uint64, payloadMessage map[string]interface{})) *ReceiverService_Notify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *ReceiverService_Notify_Call) Return(_a0 error) *ReceiverService_Notify_Call {
	_c.Call.Return(_a0)
	return _c
}

// Update provides a mock function with given fields: ctx, rcv
func (_m *ReceiverService) Update(ctx context.Context, rcv *receiver.Receiver) error {
	ret := _m.Called(ctx, rcv)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *receiver.Receiver) error); ok {
		r0 = rf(ctx, rcv)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReceiverService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ReceiverService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//  - ctx context.Context
//  - rcv *receiver.Receiver
func (_e *ReceiverService_Expecter) Update(ctx interface{}, rcv interface{}) *ReceiverService_Update_Call {
	return &ReceiverService_Update_Call{Call: _e.mock.On("Update", ctx, rcv)}
}

func (_c *ReceiverService_Update_Call) Run(run func(ctx context.Context, rcv *receiver.Receiver)) *ReceiverService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*receiver.Receiver))
	})
	return _c
}

func (_c *ReceiverService_Update_Call) Return(_a0 error) *ReceiverService_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewReceiverService interface {
	mock.TestingT
	Cleanup(func())
}

// NewReceiverService creates a new instance of ReceiverService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReceiverService(t mockConstructorTestingTNewReceiverService) *ReceiverService {
	mock := &ReceiverService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
