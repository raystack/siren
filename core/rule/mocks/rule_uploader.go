// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	provider "github.com/goto/siren/core/provider"
	mock "github.com/stretchr/testify/mock"

	rule "github.com/goto/siren/core/rule"

	template "github.com/goto/siren/core/template"
)

// RuleUploader is an autogenerated mock type for the RuleUploader type
type RuleUploader struct {
	mock.Mock
}

type RuleUploader_Expecter struct {
	mock *mock.Mock
}

func (_m *RuleUploader) EXPECT() *RuleUploader_Expecter {
	return &RuleUploader_Expecter{mock: &_m.Mock}
}

// UpsertRule provides a mock function with given fields: ctx, namespaceURN, prov, rl, templateToUpdate
func (_m *RuleUploader) UpsertRule(ctx context.Context, namespaceURN string, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template) error {
	ret := _m.Called(ctx, namespaceURN, prov, rl, templateToUpdate)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, provider.Provider, *rule.Rule, *template.Template) error); ok {
		r0 = rf(ctx, namespaceURN, prov, rl, templateToUpdate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RuleUploader_UpsertRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertRule'
type RuleUploader_UpsertRule_Call struct {
	*mock.Call
}

// UpsertRule is a helper method to define mock.On call
//   - ctx context.Context
//   - namespaceURN string
//   - prov provider.Provider
//   - rl *rule.Rule
//   - templateToUpdate *template.Template
func (_e *RuleUploader_Expecter) UpsertRule(ctx interface{}, namespaceURN interface{}, prov interface{}, rl interface{}, templateToUpdate interface{}) *RuleUploader_UpsertRule_Call {
	return &RuleUploader_UpsertRule_Call{Call: _e.mock.On("UpsertRule", ctx, namespaceURN, prov, rl, templateToUpdate)}
}

func (_c *RuleUploader_UpsertRule_Call) Run(run func(ctx context.Context, namespaceURN string, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template)) *RuleUploader_UpsertRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(provider.Provider), args[3].(*rule.Rule), args[4].(*template.Template))
	})
	return _c
}

func (_c *RuleUploader_UpsertRule_Call) Return(_a0 error) *RuleUploader_UpsertRule_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewRuleUploader interface {
	mock.TestingT
	Cleanup(func())
}

// NewRuleUploader creates a new instance of RuleUploader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuleUploader(t mockConstructorTestingTNewRuleUploader) *RuleUploader {
	mock := &RuleUploader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
