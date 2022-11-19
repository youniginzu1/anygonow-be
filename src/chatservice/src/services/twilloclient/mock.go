// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aqaurius6666/chatservice/src/services/twilloclient (interfaces: Twilio)

// Package twilloclient is a generated GoMock package.
package twilloclient

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTwilio is a mock of Twilio interface.
type MockTwilio struct {
	ctrl     *gomock.Controller
	recorder *MockTwilioMockRecorder
}

// MockTwilioMockRecorder is the mock recorder for MockTwilio.
type MockTwilioMockRecorder struct {
	mock *MockTwilio
}

// NewMockTwilio creates a new mock instance.
func NewMockTwilio(ctrl *gomock.Controller) *MockTwilio {
	mock := &MockTwilio{ctrl: ctrl}
	mock.recorder = &MockTwilioMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwilio) EXPECT() *MockTwilioMockRecorder {
	return m.recorder
}

// BuyPhoneNumber mocks base method.
func (m *MockTwilio) BuyPhoneNumber(arg0 context.Context, arg1 *string) (*string, *string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyPhoneNumber", arg0, arg1)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(*string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// BuyPhoneNumber indicates an expected call of BuyPhoneNumber.
func (mr *MockTwilioMockRecorder) BuyPhoneNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyPhoneNumber", reflect.TypeOf((*MockTwilio)(nil).BuyPhoneNumber), arg0, arg1)
}

// DeleteConversation mocks base method.
func (m *MockTwilio) DeleteConversation(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConversation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConversation indicates an expected call of DeleteConversation.
func (mr *MockTwilioMockRecorder) DeleteConversation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConversation", reflect.TypeOf((*MockTwilio)(nil).DeleteConversation), arg0, arg1)
}

// ListAvailablePhoneNumber mocks base method.
func (m *MockTwilio) ListAvailablePhoneNumber(arg0 context.Context) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAvailablePhoneNumber", arg0)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAvailablePhoneNumber indicates an expected call of ListAvailablePhoneNumber.
func (mr *MockTwilioMockRecorder) ListAvailablePhoneNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAvailablePhoneNumber", reflect.TypeOf((*MockTwilio)(nil).ListAvailablePhoneNumber), arg0)
}

// ListResourcePhone mocks base method.
func (m *MockTwilio) ListResourcePhone(arg0 context.Context) ([]string, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResourcePhone", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListResourcePhone indicates an expected call of ListResourcePhone.
func (mr *MockTwilioMockRecorder) ListResourcePhone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResourcePhone", reflect.TypeOf((*MockTwilio)(nil).ListResourcePhone), arg0)
}

// NewConversation mocks base method.
func (m *MockTwilio) NewConversation(arg0 context.Context, arg1 *string, arg2 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewConversation", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewConversation indicates an expected call of NewConversation.
func (mr *MockTwilioMockRecorder) NewConversation(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConversation", reflect.TypeOf((*MockTwilio)(nil).NewConversation), varargs...)
}

// ReleasePhoneNumber mocks base method.
func (m *MockTwilio) ReleasePhoneNumber(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleasePhoneNumber", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleasePhoneNumber indicates an expected call of ReleasePhoneNumber.
func (mr *MockTwilioMockRecorder) ReleasePhoneNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleasePhoneNumber", reflect.TypeOf((*MockTwilio)(nil).ReleasePhoneNumber), arg0, arg1)
}

// SendMessage mocks base method.
func (m *MockTwilio) SendMessage(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockTwilioMockRecorder) SendMessage(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockTwilio)(nil).SendMessage), arg0, arg1, arg2, arg3)
}