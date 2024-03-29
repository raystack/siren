// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	provider "github.com/odpf/siren/core/provider"
)

// ConfigSyncer is an autogenerated mock type for the ConfigSyncer type
type ConfigSyncer struct {
	mock.Mock
}

type ConfigSyncer_Expecter struct {
	mock *mock.Mock
}

func (_m *ConfigSyncer) EXPECT() *ConfigSyncer_Expecter {
	return &ConfigSyncer_Expecter{mock: &_m.Mock}
}

// SyncRuntimeConfig provides a mock function with given fields: ctx, namespaceID, namespaceURN, prov
func (_m *ConfigSyncer) SyncRuntimeConfig(ctx context.Context, namespaceID uint64, namespaceURN string, prov provider.Provider) error {
	ret := _m.Called(ctx, namespaceID, namespaceURN, prov)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string, provider.Provider) error); ok {
		r0 = rf(ctx, namespaceID, namespaceURN, prov)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConfigSyncer_SyncRuntimeConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncRuntimeConfig'
type ConfigSyncer_SyncRuntimeConfig_Call struct {
	*mock.Call
}

// SyncRuntimeConfig is a helper method to define mock.On call
//   - ctx context.Context
//   - namespaceID uint64
//   - namespaceURN string
//   - prov provider.Provider
func (_e *ConfigSyncer_Expecter) SyncRuntimeConfig(ctx interface{}, namespaceID interface{}, namespaceURN interface{}, prov interface{}) *ConfigSyncer_SyncRuntimeConfig_Call {
	return &ConfigSyncer_SyncRuntimeConfig_Call{Call: _e.mock.On("SyncRuntimeConfig", ctx, namespaceID, namespaceURN, prov)}
}

func (_c *ConfigSyncer_SyncRuntimeConfig_Call) Run(run func(ctx context.Context, namespaceID uint64, namespaceURN string, prov provider.Provider)) *ConfigSyncer_SyncRuntimeConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(string), args[3].(provider.Provider))
	})
	return _c
}

func (_c *ConfigSyncer_SyncRuntimeConfig_Call) Return(_a0 error) *ConfigSyncer_SyncRuntimeConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewConfigSyncer interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfigSyncer creates a new instance of ConfigSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfigSyncer(t mockConstructorTestingTNewConfigSyncer) *ConfigSyncer {
	mock := &ConfigSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
