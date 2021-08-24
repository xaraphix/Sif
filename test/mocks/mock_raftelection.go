// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xaraphix/Sif/internal/raft (interfaces: RaftElection)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/xaraphix/Sif/internal/raft"
	protos "github.com/xaraphix/Sif/internal/raft/protos"
)

// MockRaftElection is a mock of RaftElection interface.
type MockRaftElection struct {
	ctrl     *gomock.Controller
	recorder *MockRaftElectionMockRecorder
}

// MockRaftElectionMockRecorder is the mock recorder for MockRaftElection.
type MockRaftElectionMockRecorder struct {
	mock *MockRaftElection
}

// NewMockRaftElection creates a new mock instance.
func NewMockRaftElection(ctrl *gomock.Controller) *MockRaftElection {
	mock := &MockRaftElection{ctrl: ctrl}
	mock.recorder = &MockRaftElectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaftElection) EXPECT() *MockRaftElectionMockRecorder {
	return m.recorder
}

// BecomeACandidate mocks base method.
func (m *MockRaftElection) BecomeACandidate(arg0 *raft.RaftNode) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BecomeACandidate", arg0)
}

// BecomeACandidate indicates an expected call of BecomeACandidate.
func (mr *MockRaftElectionMockRecorder) BecomeACandidate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BecomeACandidate", reflect.TypeOf((*MockRaftElection)(nil).BecomeACandidate), arg0)
}

// GenerateVoteRequest mocks base method.
func (m *MockRaftElection) GenerateVoteRequest(arg0 *raft.RaftNode) *protos.VoteRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateVoteRequest", arg0)
	ret0, _ := ret[0].(*protos.VoteRequest)
	return ret0
}

// GenerateVoteRequest indicates an expected call of GenerateVoteRequest.
func (mr *MockRaftElectionMockRecorder) GenerateVoteRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateVoteRequest", reflect.TypeOf((*MockRaftElection)(nil).GenerateVoteRequest), arg0)
}

// GetElectionTimeoutDuration mocks base method.
func (m *MockRaftElection) GetElectionTimeoutDuration() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetElectionTimeoutDuration")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetElectionTimeoutDuration indicates an expected call of GetElectionTimeoutDuration.
func (mr *MockRaftElectionMockRecorder) GetElectionTimeoutDuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetElectionTimeoutDuration", reflect.TypeOf((*MockRaftElection)(nil).GetElectionTimeoutDuration))
}

// GetLeaderHeartChannel mocks base method.
func (m *MockRaftElection) GetLeaderHeartChannel() chan *raft.RaftNode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLeaderHeartChannel")
	ret0, _ := ret[0].(chan *raft.RaftNode)
	return ret0
}

// GetLeaderHeartChannel indicates an expected call of GetLeaderHeartChannel.
func (mr *MockRaftElectionMockRecorder) GetLeaderHeartChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeaderHeartChannel", reflect.TypeOf((*MockRaftElection)(nil).GetLeaderHeartChannel))
}

// GetReceivedVotes mocks base method.
func (m *MockRaftElection) GetReceivedVotes() []*protos.VoteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReceivedVotes")
	ret0, _ := ret[0].([]*protos.VoteResponse)
	return ret0
}

// GetReceivedVotes indicates an expected call of GetReceivedVotes.
func (mr *MockRaftElectionMockRecorder) GetReceivedVotes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceivedVotes", reflect.TypeOf((*MockRaftElection)(nil).GetReceivedVotes))
}

// GetResponseForVoteRequest mocks base method.
func (m *MockRaftElection) GetResponseForVoteRequest(arg0 *raft.RaftNode, arg1 *protos.VoteRequest) (*protos.VoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResponseForVoteRequest", arg0, arg1)
	ret0, _ := ret[0].(*protos.VoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResponseForVoteRequest indicates an expected call of GetResponseForVoteRequest.
func (mr *MockRaftElectionMockRecorder) GetResponseForVoteRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponseForVoteRequest", reflect.TypeOf((*MockRaftElection)(nil).GetResponseForVoteRequest), arg0, arg1)
}

// ManageElection mocks base method.
func (m *MockRaftElection) ManageElection(arg0 *raft.RaftNode) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManageElection", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ManageElection indicates an expected call of ManageElection.
func (mr *MockRaftElectionMockRecorder) ManageElection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManageElection", reflect.TypeOf((*MockRaftElection)(nil).ManageElection), arg0)
}

// StartElection mocks base method.
func (m *MockRaftElection) StartElection(arg0 *raft.RaftNode) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartElection", arg0)
}

// StartElection indicates an expected call of StartElection.
func (mr *MockRaftElectionMockRecorder) StartElection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartElection", reflect.TypeOf((*MockRaftElection)(nil).StartElection), arg0)
}
