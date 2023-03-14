// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	secret "github.com/goto/siren/pkg/secret"
	mock "github.com/stretchr/testify/mock"

	slack "github.com/goto/siren/plugins/receivers/slack"
)

// SlackCaller is an autogenerated mock type for the SlackCaller type
type SlackCaller struct {
	mock.Mock
}

type SlackCaller_Expecter struct {
	mock *mock.Mock
}

func (_m *SlackCaller) EXPECT() *SlackCaller_Expecter {
	return &SlackCaller_Expecter{mock: &_m.Mock}
}

// ExchangeAuth provides a mock function with given fields: ctx, authCode, clientID, clientSecret
func (_m *SlackCaller) ExchangeAuth(ctx context.Context, authCode string, clientID string, clientSecret string) (slack.Credential, error) {
	ret := _m.Called(ctx, authCode, clientID, clientSecret)

	var r0 slack.Credential
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) slack.Credential); ok {
		r0 = rf(ctx, authCode, clientID, clientSecret)
	} else {
		r0 = ret.Get(0).(slack.Credential)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, authCode, clientID, clientSecret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SlackCaller_ExchangeAuth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExchangeAuth'
type SlackCaller_ExchangeAuth_Call struct {
	*mock.Call
}

// ExchangeAuth is a helper method to define mock.On call
//   - ctx context.Context
//   - authCode string
//   - clientID string
//   - clientSecret string
func (_e *SlackCaller_Expecter) ExchangeAuth(ctx interface{}, authCode interface{}, clientID interface{}, clientSecret interface{}) *SlackCaller_ExchangeAuth_Call {
	return &SlackCaller_ExchangeAuth_Call{Call: _e.mock.On("ExchangeAuth", ctx, authCode, clientID, clientSecret)}
}

func (_c *SlackCaller_ExchangeAuth_Call) Run(run func(ctx context.Context, authCode string, clientID string, clientSecret string)) *SlackCaller_ExchangeAuth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *SlackCaller_ExchangeAuth_Call) Return(_a0 slack.Credential, _a1 error) *SlackCaller_ExchangeAuth_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetWorkspaceChannels provides a mock function with given fields: ctx, token
func (_m *SlackCaller) GetWorkspaceChannels(ctx context.Context, token secret.MaskableString) ([]slack.Channel, error) {
	ret := _m.Called(ctx, token)

	var r0 []slack.Channel
	if rf, ok := ret.Get(0).(func(context.Context, secret.MaskableString) []slack.Channel); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]slack.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, secret.MaskableString) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SlackCaller_GetWorkspaceChannels_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWorkspaceChannels'
type SlackCaller_GetWorkspaceChannels_Call struct {
	*mock.Call
}

// GetWorkspaceChannels is a helper method to define mock.On call
//   - ctx context.Context
//   - token secret.MaskableString
func (_e *SlackCaller_Expecter) GetWorkspaceChannels(ctx interface{}, token interface{}) *SlackCaller_GetWorkspaceChannels_Call {
	return &SlackCaller_GetWorkspaceChannels_Call{Call: _e.mock.On("GetWorkspaceChannels", ctx, token)}
}

func (_c *SlackCaller_GetWorkspaceChannels_Call) Run(run func(ctx context.Context, token secret.MaskableString)) *SlackCaller_GetWorkspaceChannels_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(secret.MaskableString))
	})
	return _c
}

func (_c *SlackCaller_GetWorkspaceChannels_Call) Return(_a0 []slack.Channel, _a1 error) *SlackCaller_GetWorkspaceChannels_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Notify provides a mock function with given fields: ctx, conf, message
func (_m *SlackCaller) Notify(ctx context.Context, conf slack.NotificationConfig, message slack.Message) error {
	ret := _m.Called(ctx, conf, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, slack.NotificationConfig, slack.Message) error); ok {
		r0 = rf(ctx, conf, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SlackCaller_Notify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Notify'
type SlackCaller_Notify_Call struct {
	*mock.Call
}

// Notify is a helper method to define mock.On call
//   - ctx context.Context
//   - conf slack.NotificationConfig
//   - message slack.Message
func (_e *SlackCaller_Expecter) Notify(ctx interface{}, conf interface{}, message interface{}) *SlackCaller_Notify_Call {
	return &SlackCaller_Notify_Call{Call: _e.mock.On("Notify", ctx, conf, message)}
}

func (_c *SlackCaller_Notify_Call) Run(run func(ctx context.Context, conf slack.NotificationConfig, message slack.Message)) *SlackCaller_Notify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(slack.NotificationConfig), args[2].(slack.Message))
	})
	return _c
}

func (_c *SlackCaller_Notify_Call) Return(_a0 error) *SlackCaller_Notify_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewSlackCaller interface {
	mock.TestingT
	Cleanup(func())
}

// NewSlackCaller creates a new instance of SlackCaller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSlackCaller(t mockConstructorTestingTNewSlackCaller) *SlackCaller {
	mock := &SlackCaller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
