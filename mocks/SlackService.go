// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	slack "github.com/slack-go/slack"
	mock "github.com/stretchr/testify/mock"
)

// SlackService is an autogenerated mock type for the SlackService type
type SlackService struct {
	mock.Mock
}

// GetJoinedChannelsList provides a mock function with given fields:
func (_m *SlackService) GetJoinedChannelsList() ([]slack.Channel, error) {
	ret := _m.Called()

	var r0 []slack.Channel
	if rf, ok := ret.Get(0).(func() []slack.Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]slack.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: _a0
func (_m *SlackService) GetUserByEmail(_a0 string) (*slack.User, error) {
	ret := _m.Called(_a0)

	var r0 *slack.User
	if rf, ok := ret.Get(0).(func(string) *slack.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*slack.User)
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

// SendMessage provides a mock function with given fields: _a0, _a1
func (_m *SlackService) SendMessage(_a0 string, _a1 ...slack.MsgOption) (string, string, string, error) {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, ...slack.MsgOption) string); ok {
		r0 = rf(_a0, _a1...)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, ...slack.MsgOption) string); ok {
		r1 = rf(_a0, _a1...)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(string, ...slack.MsgOption) string); ok {
		r2 = rf(_a0, _a1...)
	} else {
		r2 = ret.Get(2).(string)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string, ...slack.MsgOption) error); ok {
		r3 = rf(_a0, _a1...)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

