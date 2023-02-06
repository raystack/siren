// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	rule "github.com/odpf/siren/core/rule"
	mock "github.com/stretchr/testify/mock"
)

// RuleRepository is an autogenerated mock type for the Repository type
type RuleRepository struct {
	mock.Mock
}

type RuleRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *RuleRepository) EXPECT() *RuleRepository_Expecter {
	return &RuleRepository_Expecter{mock: &_m.Mock}
}

// Commit provides a mock function with given fields: ctx
func (_m *RuleRepository) Commit(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RuleRepository_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type RuleRepository_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RuleRepository_Expecter) Commit(ctx interface{}) *RuleRepository_Commit_Call {
	return &RuleRepository_Commit_Call{Call: _e.mock.On("Commit", ctx)}
}

func (_c *RuleRepository_Commit_Call) Run(run func(ctx context.Context)) *RuleRepository_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RuleRepository_Commit_Call) Return(_a0 error) *RuleRepository_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

// List provides a mock function with given fields: _a0, _a1
func (_m *RuleRepository) List(_a0 context.Context, _a1 rule.Filter) ([]rule.Rule, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []rule.Rule
	if rf, ok := ret.Get(0).(func(context.Context, rule.Filter) []rule.Rule); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]rule.Rule)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, rule.Filter) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RuleRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type RuleRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 rule.Filter
func (_e *RuleRepository_Expecter) List(_a0 interface{}, _a1 interface{}) *RuleRepository_List_Call {
	return &RuleRepository_List_Call{Call: _e.mock.On("List", _a0, _a1)}
}

func (_c *RuleRepository_List_Call) Run(run func(_a0 context.Context, _a1 rule.Filter)) *RuleRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(rule.Filter))
	})
	return _c
}

func (_c *RuleRepository_List_Call) Return(_a0 []rule.Rule, _a1 error) *RuleRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Rollback provides a mock function with given fields: ctx, err
func (_m *RuleRepository) Rollback(ctx context.Context, err error) error {
	ret := _m.Called(ctx, err)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, error) error); ok {
		r0 = rf(ctx, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RuleRepository_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type RuleRepository_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *RuleRepository_Expecter) Rollback(ctx interface{}, err interface{}) *RuleRepository_Rollback_Call {
	return &RuleRepository_Rollback_Call{Call: _e.mock.On("Rollback", ctx, err)}
}

func (_c *RuleRepository_Rollback_Call) Run(run func(ctx context.Context, err error)) *RuleRepository_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *RuleRepository_Rollback_Call) Return(_a0 error) *RuleRepository_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

// Upsert provides a mock function with given fields: _a0, _a1
func (_m *RuleRepository) Upsert(_a0 context.Context, _a1 *rule.Rule) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *rule.Rule) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RuleRepository_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type RuleRepository_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *rule.Rule
func (_e *RuleRepository_Expecter) Upsert(_a0 interface{}, _a1 interface{}) *RuleRepository_Upsert_Call {
	return &RuleRepository_Upsert_Call{Call: _e.mock.On("Upsert", _a0, _a1)}
}

func (_c *RuleRepository_Upsert_Call) Run(run func(_a0 context.Context, _a1 *rule.Rule)) *RuleRepository_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*rule.Rule))
	})
	return _c
}

func (_c *RuleRepository_Upsert_Call) Return(_a0 error) *RuleRepository_Upsert_Call {
	_c.Call.Return(_a0)
	return _c
}

// WithTransaction provides a mock function with given fields: ctx
func (_m *RuleRepository) WithTransaction(ctx context.Context) context.Context {
	ret := _m.Called(ctx)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(context.Context) context.Context); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// RuleRepository_WithTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTransaction'
type RuleRepository_WithTransaction_Call struct {
	*mock.Call
}

// WithTransaction is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RuleRepository_Expecter) WithTransaction(ctx interface{}) *RuleRepository_WithTransaction_Call {
	return &RuleRepository_WithTransaction_Call{Call: _e.mock.On("WithTransaction", ctx)}
}

func (_c *RuleRepository_WithTransaction_Call) Run(run func(ctx context.Context)) *RuleRepository_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RuleRepository_WithTransaction_Call) Return(_a0 context.Context) *RuleRepository_WithTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewRuleRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRuleRepository creates a new instance of RuleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuleRepository(t mockConstructorTestingTNewRuleRepository) *RuleRepository {
	mock := &RuleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
