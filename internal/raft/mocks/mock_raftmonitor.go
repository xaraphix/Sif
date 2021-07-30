// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xaraphix/Sif/internal/raft (interfaces: RaftMonitor)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/xaraphix/Sif/internal/raft"
)

// MockRaftMonitor is a mock of RaftMonitor interface.
type MockRaftMonitor struct {
	ctrl     *gomock.Controller
	recorder *MockRaftMonitorMockRecorder
}

// MockRaftMonitorMockRecorder is the mock recorder for MockRaftMonitor.
type MockRaftMonitorMockRecorder struct {
	mock *MockRaftMonitor
}

// NewMockRaftMonitor creates a new mock instance.
func NewMockRaftMonitor(ctrl *gomock.Controller) *MockRaftMonitor {
	mock := &MockRaftMonitor{ctrl: ctrl}
	mock.recorder = &MockRaftMonitorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaftMonitor) EXPECT() *MockRaftMonitorMockRecorder {
	return m.recorder
}

// GetLastResetAt mocks base method.
func (m *MockRaftMonitor) GetLastResetAt() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastResetAt")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetLastResetAt indicates an expected call of GetLastResetAt.
func (mr *MockRaftMonitorMockRecorder) GetLastResetAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastResetAt", reflect.TypeOf((*MockRaftMonitor)(nil).GetLastResetAt))
}

// IsAutoStartOn mocks base method.
func (m *MockRaftMonitor) IsAutoStartOn() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAutoStartOn")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAutoStartOn indicates an expected call of IsAutoStartOn.
func (mr *MockRaftMonitorMockRecorder) IsAutoStartOn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAutoStartOn", reflect.TypeOf((*MockRaftMonitor)(nil).IsAutoStartOn))
}

// Sleep mocks base method.
func (m *MockRaftMonitor) Sleep() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Sleep")
}

// Sleep indicates an expected call of Sleep.
func (mr *MockRaftMonitorMockRecorder) Sleep() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sleep", reflect.TypeOf((*MockRaftMonitor)(nil).Sleep))
}

// Start mocks base method.
func (m *MockRaftMonitor) Start(arg0 *raft.RaftNode) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0)
}

// Start indicates an expected call of Start.
func (mr *MockRaftMonitorMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRaftMonitor)(nil).Start), arg0)
}

// Stop mocks base method.
func (m *MockRaftMonitor) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockRaftMonitorMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockRaftMonitor)(nil).Stop))
}
