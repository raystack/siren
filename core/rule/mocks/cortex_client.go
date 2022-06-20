// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	cortex "github.com/odpf/siren/pkg/cortex"
	mock "github.com/stretchr/testify/mock"

	rwrulefmt "github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"

	testing "testing"
)

// CortexClient is an autogenerated mock type for the CortexClient type
type CortexClient struct {
	mock.Mock
}

type CortexClient_Expecter struct {
	mock *mock.Mock
}

func (_m *CortexClient) EXPECT() *CortexClient_Expecter {
	return &CortexClient_Expecter{mock: &_m.Mock}
}

// CreateAlertmanagerConfig provides a mock function with given fields: _a0, _a1
func (_m *CortexClient) CreateAlertmanagerConfig(_a0 cortex.AlertManagerConfig, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(cortex.AlertManagerConfig, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CortexClient_CreateAlertmanagerConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAlertmanagerConfig'
type CortexClient_CreateAlertmanagerConfig_Call struct {
	*mock.Call
}

// CreateAlertmanagerConfig is a helper method to define mock.On call
//  - _a0 cortex.AlertManagerConfig
//  - _a1 string
func (_e *CortexClient_Expecter) CreateAlertmanagerConfig(_a0 interface{}, _a1 interface{}) *CortexClient_CreateAlertmanagerConfig_Call {
	return &CortexClient_CreateAlertmanagerConfig_Call{Call: _e.mock.On("CreateAlertmanagerConfig", _a0, _a1)}
}

func (_c *CortexClient_CreateAlertmanagerConfig_Call) Run(run func(_a0 cortex.AlertManagerConfig, _a1 string)) *CortexClient_CreateAlertmanagerConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(cortex.AlertManagerConfig), args[1].(string))
	})
	return _c
}

func (_c *CortexClient_CreateAlertmanagerConfig_Call) Return(_a0 error) *CortexClient_CreateAlertmanagerConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

// CreateRuleGroup provides a mock function with given fields: ctx, namespace, rg
func (_m *CortexClient) CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error {
	ret := _m.Called(ctx, namespace, rg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, rwrulefmt.RuleGroup) error); ok {
		r0 = rf(ctx, namespace, rg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CortexClient_CreateRuleGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRuleGroup'
type CortexClient_CreateRuleGroup_Call struct {
	*mock.Call
}

// CreateRuleGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - namespace string
//  - rg rwrulefmt.RuleGroup
func (_e *CortexClient_Expecter) CreateRuleGroup(ctx interface{}, namespace interface{}, rg interface{}) *CortexClient_CreateRuleGroup_Call {
	return &CortexClient_CreateRuleGroup_Call{Call: _e.mock.On("CreateRuleGroup", ctx, namespace, rg)}
}

func (_c *CortexClient_CreateRuleGroup_Call) Run(run func(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup)) *CortexClient_CreateRuleGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(rwrulefmt.RuleGroup))
	})
	return _c
}

func (_c *CortexClient_CreateRuleGroup_Call) Return(_a0 error) *CortexClient_CreateRuleGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

// DeleteRuleGroup provides a mock function with given fields: ctx, namespace, groupName
func (_m *CortexClient) DeleteRuleGroup(ctx context.Context, namespace string, groupName string) error {
	ret := _m.Called(ctx, namespace, groupName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, namespace, groupName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CortexClient_DeleteRuleGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRuleGroup'
type CortexClient_DeleteRuleGroup_Call struct {
	*mock.Call
}

// DeleteRuleGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - namespace string
//  - groupName string
func (_e *CortexClient_Expecter) DeleteRuleGroup(ctx interface{}, namespace interface{}, groupName interface{}) *CortexClient_DeleteRuleGroup_Call {
	return &CortexClient_DeleteRuleGroup_Call{Call: _e.mock.On("DeleteRuleGroup", ctx, namespace, groupName)}
}

func (_c *CortexClient_DeleteRuleGroup_Call) Run(run func(ctx context.Context, namespace string, groupName string)) *CortexClient_DeleteRuleGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CortexClient_DeleteRuleGroup_Call) Return(_a0 error) *CortexClient_DeleteRuleGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetRuleGroup provides a mock function with given fields: ctx, namespace, groupName
func (_m *CortexClient) GetRuleGroup(ctx context.Context, namespace string, groupName string) (*rwrulefmt.RuleGroup, error) {
	ret := _m.Called(ctx, namespace, groupName)

	var r0 *rwrulefmt.RuleGroup
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *rwrulefmt.RuleGroup); ok {
		r0 = rf(ctx, namespace, groupName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rwrulefmt.RuleGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, namespace, groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CortexClient_GetRuleGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRuleGroup'
type CortexClient_GetRuleGroup_Call struct {
	*mock.Call
}

// GetRuleGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - namespace string
//  - groupName string
func (_e *CortexClient_Expecter) GetRuleGroup(ctx interface{}, namespace interface{}, groupName interface{}) *CortexClient_GetRuleGroup_Call {
	return &CortexClient_GetRuleGroup_Call{Call: _e.mock.On("GetRuleGroup", ctx, namespace, groupName)}
}

func (_c *CortexClient_GetRuleGroup_Call) Run(run func(ctx context.Context, namespace string, groupName string)) *CortexClient_GetRuleGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CortexClient_GetRuleGroup_Call) Return(_a0 *rwrulefmt.RuleGroup, _a1 error) *CortexClient_GetRuleGroup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// ListRules provides a mock function with given fields: ctx, namespace
func (_m *CortexClient) ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error) {
	ret := _m.Called(ctx, namespace)

	var r0 map[string][]rwrulefmt.RuleGroup
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string][]rwrulefmt.RuleGroup); ok {
		r0 = rf(ctx, namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]rwrulefmt.RuleGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CortexClient_ListRules_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListRules'
type CortexClient_ListRules_Call struct {
	*mock.Call
}

// ListRules is a helper method to define mock.On call
//  - ctx context.Context
//  - namespace string
func (_e *CortexClient_Expecter) ListRules(ctx interface{}, namespace interface{}) *CortexClient_ListRules_Call {
	return &CortexClient_ListRules_Call{Call: _e.mock.On("ListRules", ctx, namespace)}
}

func (_c *CortexClient_ListRules_Call) Run(run func(ctx context.Context, namespace string)) *CortexClient_ListRules_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CortexClient_ListRules_Call) Return(_a0 map[string][]rwrulefmt.RuleGroup, _a1 error) *CortexClient_ListRules_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// NewCortexClient creates a new instance of CortexClient. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewCortexClient(t testing.TB) *CortexClient {
	mock := &CortexClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
