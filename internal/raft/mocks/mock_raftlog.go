// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xaraphix/Sif/internal/raft (interfaces: RaftLog)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/xaraphix/Sif/internal/raft"
)

// MockRaftLog is a mock of RaftLog interface.
type MockRaftLog struct {
	ctrl     *gomock.Controller
	recorder *MockRaftLogMockRecorder
}

// MockRaftLogMockRecorder is the mock recorder for MockRaftLog.
type MockRaftLogMockRecorder struct {
	mock *MockRaftLog
}

// NewMockRaftLog creates a new mock instance.
func NewMockRaftLog(ctrl *gomock.Controller) *MockRaftLog {
	mock := &MockRaftLog{ctrl: ctrl}
	mock.recorder = &MockRaftLogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaftLog) EXPECT() *MockRaftLogMockRecorder {
	return m.recorder
}

// GetLog mocks base method.
func (m *MockRaftLog) GetLog(arg0 int32) raft.Log {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLog", arg0)
	ret0, _ := ret[0].(raft.Log)
	return ret0
}

// GetLog indicates an expected call of GetLog.
func (mr *MockRaftLogMockRecorder) GetLog(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLog", reflect.TypeOf((*MockRaftLog)(nil).GetLog), arg0)
}

// GetLogs mocks base method.
func (m *MockRaftLog) GetLogs() []raft.Log {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs")
	ret0, _ := ret[0].([]raft.Log)
	return ret0
}

// GetLogs indicates an expected call of GetLogs.
func (mr *MockRaftLogMockRecorder) GetLogs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockRaftLog)(nil).GetLogs))
}