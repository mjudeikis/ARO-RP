// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azure/ARO-RP/pkg/operator/controllers/workaround (interfaces: Workaround)

// Package mock_workaround is a generated GoMock package.
package mock_workaround

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	version "github.com/Azure/ARO-RP/pkg/util/version"
)

// MockWorkaround is a mock of Workaround interface
type MockWorkaround struct {
	ctrl     *gomock.Controller
	recorder *MockWorkaroundMockRecorder
}

// MockWorkaroundMockRecorder is the mock recorder for MockWorkaround
type MockWorkaroundMockRecorder struct {
	mock *MockWorkaround
}

// NewMockWorkaround creates a new mock instance
func NewMockWorkaround(ctrl *gomock.Controller) *MockWorkaround {
	mock := &MockWorkaround{ctrl: ctrl}
	mock.recorder = &MockWorkaroundMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkaround) EXPECT() *MockWorkaroundMockRecorder {
	return m.recorder
}

// Ensure mocks base method
func (m *MockWorkaround) Ensure(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ensure", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ensure indicates an expected call of Ensure
func (mr *MockWorkaroundMockRecorder) Ensure(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ensure", reflect.TypeOf((*MockWorkaround)(nil).Ensure), arg0)
}

// IsRequired mocks base method
func (m *MockWorkaround) IsRequired(arg0 *version.Version) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRequired", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRequired indicates an expected call of IsRequired
func (mr *MockWorkaroundMockRecorder) IsRequired(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRequired", reflect.TypeOf((*MockWorkaround)(nil).IsRequired), arg0)
}

// Name mocks base method
func (m *MockWorkaround) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockWorkaroundMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockWorkaround)(nil).Name))
}

// Remove mocks base method
func (m *MockWorkaround) Remove(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockWorkaroundMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockWorkaround)(nil).Remove), arg0)
}
