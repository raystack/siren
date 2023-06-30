// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	provider "github.com/raystack/siren/core/provider"
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

// Get provides a mock function with given fields: ctx, id
func (_m *ProviderService) Get(ctx context.Context, id uint64) (*provider.Provider, error) {
	ret := _m.Called(ctx, id)

	var r0 *provider.Provider
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *provider.Provider); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*provider.Provider)
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

// ProviderService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ProviderService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint64
func (_e *ProviderService_Expecter) Get(ctx interface{}, id interface{}) *ProviderService_Get_Call {
	return &ProviderService_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *ProviderService_Get_Call) Run(run func(ctx context.Context, id uint64)) *ProviderService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *ProviderService_Get_Call) Return(_a0 *provider.Provider, _a1 error) *ProviderService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewProviderService interface {
	mock.TestingT
	Cleanup(func())
}

// NewProviderService creates a new instance of ProviderService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProviderService(t mockConstructorTestingTNewProviderService) *ProviderService {
	mock := &ProviderService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
