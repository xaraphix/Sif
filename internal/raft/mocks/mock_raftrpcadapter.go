// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xaraphix/Sif/internal/raft (interfaces: RaftRPCAdapter)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/xaraphix/Sif/internal/raft"
)

// MockRaftRPCAdapter is a mock of RaftRPCAdapter interface.
type MockRaftRPCAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockRaftRPCAdapterMockRecorder
}

// MockRaftRPCAdapterMockRecorder is the mock recorder for MockRaftRPCAdapter.
type MockRaftRPCAdapterMockRecorder struct {
	mock *MockRaftRPCAdapter
}

// NewMockRaftRPCAdapter creates a new mock instance.
func NewMockRaftRPCAdapter(ctrl *gomock.Controller) *MockRaftRPCAdapter {
	mock := &MockRaftRPCAdapter{ctrl: ctrl}
	mock.recorder = &MockRaftRPCAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaftRPCAdapter) EXPECT() *MockRaftRPCAdapterMockRecorder {
	return m.recorder
}

// GenerateVoteResponse mocks base method.
func (m *MockRaftRPCAdapter) GenerateVoteResponse(arg0 raft.VoteRequest) raft.VoteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateVoteResponse", arg0)
	ret0, _ := ret[0].(raft.VoteResponse)
	return ret0
}

// GenerateVoteResponse indicates an expected call of GenerateVoteResponse.
func (mr *MockRaftRPCAdapterMockRecorder) GenerateVoteResponse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateVoteResponse", reflect.TypeOf((*MockRaftRPCAdapter)(nil).GenerateVoteResponse), arg0)
}

// ReceiveLogRequest mocks base method.
func (m *MockRaftRPCAdapter) ReceiveLogRequest(arg0 raft.LogRequest) raft.LogResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveLogRequest", arg0)
	ret0, _ := ret[0].(raft.LogResponse)
	return ret0
}

// ReceiveLogRequest indicates an expected call of ReceiveLogRequest.
func (mr *MockRaftRPCAdapterMockRecorder) ReceiveLogRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveLogRequest", reflect.TypeOf((*MockRaftRPCAdapter)(nil).ReceiveLogRequest), arg0)
}

// ReplicateLog mocks base method.
func (m *MockRaftRPCAdapter) ReplicateLog(arg0 raft.Peer, arg1 raft.LogRequest) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReplicateLog", arg0, arg1)
}

// ReplicateLog indicates an expected call of ReplicateLog.
func (mr *MockRaftRPCAdapterMockRecorder) ReplicateLog(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateLog", reflect.TypeOf((*MockRaftRPCAdapter)(nil).ReplicateLog), arg0, arg1)
}

// RequestVoteFromPeer mocks base method.
func (m *MockRaftRPCAdapter) RequestVoteFromPeer(arg0 raft.Peer, arg1 raft.VoteRequest) raft.VoteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestVoteFromPeer", arg0, arg1)
	ret0, _ := ret[0].(raft.VoteResponse)
	return ret0
}

// RequestVoteFromPeer indicates an expected call of RequestVoteFromPeer.
func (mr *MockRaftRPCAdapterMockRecorder) RequestVoteFromPeer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestVoteFromPeer", reflect.TypeOf((*MockRaftRPCAdapter)(nil).RequestVoteFromPeer), arg0, arg1)
}
