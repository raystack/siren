// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/odpf/siren/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TemplatesService is an autogenerated mock type for the TemplatesService type
type TemplatesService struct {
	mock.Mock
}

type TemplatesService_Expecter struct {
	mock *mock.Mock
}

func (_m *TemplatesService) EXPECT() *TemplatesService_Expecter {
	return &TemplatesService_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0
func (_m *TemplatesService) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TemplatesService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type TemplatesService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//  - _a0 string
func (_e *TemplatesService_Expecter) Delete(_a0 interface{}) *TemplatesService_Delete_Call {
	return &TemplatesService_Delete_Call{Call: _e.mock.On("Delete", _a0)}
}

func (_c *TemplatesService_Delete_Call) Run(run func(_a0 string)) *TemplatesService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TemplatesService_Delete_Call) Return(_a0 error) *TemplatesService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetByName provides a mock function with given fields: _a0
func (_m *TemplatesService) GetByName(_a0 string) (*domain.Template, error) {
	ret := _m.Called(_a0)

	var r0 *domain.Template
	if rf, ok := ret.Get(0).(func(string) *domain.Template); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplatesService_GetByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByName'
type TemplatesService_GetByName_Call struct {
	*mock.Call
}

// GetByName is a helper method to define mock.On call
//  - _a0 string
func (_e *TemplatesService_Expecter) GetByName(_a0 interface{}) *TemplatesService_GetByName_Call {
	return &TemplatesService_GetByName_Call{Call: _e.mock.On("GetByName", _a0)}
}

func (_c *TemplatesService_GetByName_Call) Run(run func(_a0 string)) *TemplatesService_GetByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TemplatesService_GetByName_Call) Return(_a0 *domain.Template, _a1 error) *TemplatesService_GetByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Index provides a mock function with given fields: _a0
func (_m *TemplatesService) Index(_a0 string) ([]domain.Template, error) {
	ret := _m.Called(_a0)

	var r0 []domain.Template
	if rf, ok := ret.Get(0).(func(string) []domain.Template); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplatesService_Index_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Index'
type TemplatesService_Index_Call struct {
	*mock.Call
}

// Index is a helper method to define mock.On call
//  - _a0 string
func (_e *TemplatesService_Expecter) Index(_a0 interface{}) *TemplatesService_Index_Call {
	return &TemplatesService_Index_Call{Call: _e.mock.On("Index", _a0)}
}

func (_c *TemplatesService_Index_Call) Run(run func(_a0 string)) *TemplatesService_Index_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TemplatesService_Index_Call) Return(_a0 []domain.Template, _a1 error) *TemplatesService_Index_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Migrate provides a mock function with given fields:
func (_m *TemplatesService) Migrate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TemplatesService_Migrate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Migrate'
type TemplatesService_Migrate_Call struct {
	*mock.Call
}

// Migrate is a helper method to define mock.On call
func (_e *TemplatesService_Expecter) Migrate() *TemplatesService_Migrate_Call {
	return &TemplatesService_Migrate_Call{Call: _e.mock.On("Migrate")}
}

func (_c *TemplatesService_Migrate_Call) Run(run func()) *TemplatesService_Migrate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TemplatesService_Migrate_Call) Return(_a0 error) *TemplatesService_Migrate_Call {
	_c.Call.Return(_a0)
	return _c
}

// Render provides a mock function with given fields: _a0, _a1
func (_m *TemplatesService) Render(_a0 string, _a1 map[string]string) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, map[string]string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplatesService_Render_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Render'
type TemplatesService_Render_Call struct {
	*mock.Call
}

// Render is a helper method to define mock.On call
//  - _a0 string
//  - _a1 map[string]string
func (_e *TemplatesService_Expecter) Render(_a0 interface{}, _a1 interface{}) *TemplatesService_Render_Call {
	return &TemplatesService_Render_Call{Call: _e.mock.On("Render", _a0, _a1)}
}

func (_c *TemplatesService_Render_Call) Run(run func(_a0 string, _a1 map[string]string)) *TemplatesService_Render_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(map[string]string))
	})
	return _c
}

func (_c *TemplatesService_Render_Call) Return(_a0 string, _a1 error) *TemplatesService_Render_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Upsert provides a mock function with given fields: _a0
func (_m *TemplatesService) Upsert(_a0 *domain.Template) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Template) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TemplatesService_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type TemplatesService_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//  - _a0 *domain.Template
func (_e *TemplatesService_Expecter) Upsert(_a0 interface{}) *TemplatesService_Upsert_Call {
	return &TemplatesService_Upsert_Call{Call: _e.mock.On("Upsert", _a0)}
}

func (_c *TemplatesService_Upsert_Call) Run(run func(_a0 *domain.Template)) *TemplatesService_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Template))
	})
	return _c
}

func (_c *TemplatesService_Upsert_Call) Return(_a0 error) *TemplatesService_Upsert_Call {
	_c.Call.Return(_a0)
	return _c
}

// NewTemplatesService creates a new instance of TemplatesService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTemplatesService(t testing.TB) *TemplatesService {
	mock := &TemplatesService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
