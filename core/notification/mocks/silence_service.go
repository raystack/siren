// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	silence "github.com/raystack/siren/core/silence"
)

// SilenceService is an autogenerated mock type for the SilenceService type
type SilenceService struct {
	mock.Mock
}

type SilenceService_Expecter struct {
	mock *mock.Mock
}

func (_m *SilenceService) EXPECT() *SilenceService_Expecter {
	return &SilenceService_Expecter{mock: &_m.Mock}
}

// List provides a mock function with given fields: ctx, filter
func (_m *SilenceService) List(ctx context.Context, filter silence.Filter) ([]silence.Silence, error) {
	ret := _m.Called(ctx, filter)

	var r0 []silence.Silence
	if rf, ok := ret.Get(0).(func(context.Context, silence.Filter) []silence.Silence); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]silence.Silence)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, silence.Filter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SilenceService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type SilenceService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - filter silence.Filter
func (_e *SilenceService_Expecter) List(ctx interface{}, filter interface{}) *SilenceService_List_Call {
	return &SilenceService_List_Call{Call: _e.mock.On("List", ctx, filter)}
}

func (_c *SilenceService_List_Call) Run(run func(ctx context.Context, filter silence.Filter)) *SilenceService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(silence.Filter))
	})
	return _c
}

func (_c *SilenceService_List_Call) Return(_a0 []silence.Silence, _a1 error) *SilenceService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewSilenceService interface {
	mock.TestingT
	Cleanup(func())
}

// NewSilenceService creates a new instance of SilenceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSilenceService(t mockConstructorTestingTNewSilenceService) *SilenceService {
	mock := &SilenceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
