// Code generated by mockery 2.9.0. DO NOT EDIT.

package alerts

import mock "github.com/stretchr/testify/mock"

// MockAlertRepository is an autogenerated mock type for the AlertRepository type
type MockAlertRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *MockAlertRepository) Create(_a0 *Alert) (*Alert, error) {
	ret := _m.Called(_a0)

	var r0 *Alert
	if rf, ok := ret.Get(0).(func(*Alert) *Alert); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Alert)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Alert) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *MockAlertRepository) Get(_a0 string, _a1 uint64, _a2 uint64, _a3 uint64) ([]Alert, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []Alert
	if rf, ok := ret.Get(0).(func(string, uint64, uint64, uint64) []Alert); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Alert)
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

// Migrate provides a mock function with given fields:
func (_m *MockAlertRepository) Migrate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}