// Code generated by mockery 2.9.4. DO NOT EDIT.

package namespace

import mock "github.com/stretchr/testify/mock"

// EncryptorDecryptor is an autogenerated mock type for the EncryptorDecryptor type
type EncryptorDecryptorMock struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: _a0
func (_m *EncryptorDecryptorMock) Decrypt(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encrypt provides a mock function with given fields: _a0
func (_m *EncryptorDecryptorMock) Encrypt(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
