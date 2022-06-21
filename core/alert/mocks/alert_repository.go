// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	alert "github.com/odpf/siren/core/alert"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// AlertRepository is an autogenerated mock type for the Repository type
type AlertRepository struct {
	mock.Mock
}

type AlertRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *AlertRepository) EXPECT() *AlertRepository_Expecter {
	return &AlertRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0
func (_m *AlertRepository) Create(_a0 *alert.Alert) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*alert.Alert) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AlertRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type AlertRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//  - _a0 *alert.Alert
func (_e *AlertRepository_Expecter) Create(_a0 interface{}) *AlertRepository_Create_Call {
	return &AlertRepository_Create_Call{Call: _e.mock.On("Create", _a0)}
}

func (_c *AlertRepository_Create_Call) Run(run func(_a0 *alert.Alert)) *AlertRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*alert.Alert))
	})
	return _c
}

func (_c *AlertRepository_Create_Call) Return(_a0 error) *AlertRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) Get(_a0 string, _a1 uint64, _a2 uint64, _a3 uint64) ([]alert.Alert, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []alert.Alert
	if rf, ok := ret.Get(0).(func(string, uint64, uint64, uint64) []alert.Alert); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]alert.Alert)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, uint64, uint64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AlertRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type AlertRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//  - _a0 string
//  - _a1 uint64
//  - _a2 uint64
//  - _a3 uint64
func (_e *AlertRepository_Expecter) Get(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *AlertRepository_Get_Call {
	return &AlertRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1, _a2, _a3)}
}

func (_c *AlertRepository_Get_Call) Run(run func(_a0 string, _a1 uint64, _a2 uint64, _a3 uint64)) *AlertRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(uint64), args[2].(uint64), args[3].(uint64))
	})
	return _c
}

func (_c *AlertRepository_Get_Call) Return(_a0 []alert.Alert, _a1 error) *AlertRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// NewAlertRepository creates a new instance of AlertRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAlertRepository(t testing.TB) *AlertRepository {
	mock := &AlertRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
