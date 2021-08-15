// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xaraphix/Sif/internal/raft (interfaces: RaftConfig)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/xaraphix/Sif/internal/raft"
	protos "github.com/xaraphix/Sif/internal/raft/protos"
)

// MockRaftConfig is a mock of RaftConfig interface.
type MockRaftConfig struct {
	ctrl     *gomock.Controller
	recorder *MockRaftConfigMockRecorder
}

// MockRaftConfigMockRecorder is the mock recorder for MockRaftConfig.
type MockRaftConfigMockRecorder struct {
	mock *MockRaftConfig
}

// NewMockRaftConfig creates a new mock instance.
func NewMockRaftConfig(ctrl *gomock.Controller) *MockRaftConfig {
	mock := &MockRaftConfig{ctrl: ctrl}
	mock.recorder = &MockRaftConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaftConfig) EXPECT() *MockRaftConfigMockRecorder {
	return m.recorder
}

// CommitLength mocks base method.
func (m *MockRaftConfig) CommitLength() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitLength")
	ret0, _ := ret[0].(int32)
	return ret0
}

// CommitLength indicates an expected call of CommitLength.
func (mr *MockRaftConfigMockRecorder) CommitLength() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitLength", reflect.TypeOf((*MockRaftConfig)(nil).CommitLength))
}

// CurrentTerm mocks base method.
func (m *MockRaftConfig) CurrentTerm() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentTerm")
	ret0, _ := ret[0].(int32)
	return ret0
}

// CurrentTerm indicates an expected call of CurrentTerm.
func (mr *MockRaftConfigMockRecorder) CurrentTerm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentTerm", reflect.TypeOf((*MockRaftConfig)(nil).CurrentTerm))
}

// DidNodeCrash mocks base method.
func (m *MockRaftConfig) DidNodeCrash(arg0 *raft.RaftNode) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DidNodeCrash", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// DidNodeCrash indicates an expected call of DidNodeCrash.
func (mr *MockRaftConfigMockRecorder) DidNodeCrash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DidNodeCrash", reflect.TypeOf((*MockRaftConfig)(nil).DidNodeCrash), arg0)
}

// Host mocks base method.
func (m *MockRaftConfig) Host() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Host")
	ret0, _ := ret[0].(string)
	return ret0
}

// Host indicates an expected call of Host.
func (mr *MockRaftConfigMockRecorder) Host() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Host", reflect.TypeOf((*MockRaftConfig)(nil).Host))
}

// InstanceDirPath mocks base method.
func (m *MockRaftConfig) InstanceDirPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceDirPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// InstanceDirPath indicates an expected call of InstanceDirPath.
func (mr *MockRaftConfigMockRecorder) InstanceDirPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceDirPath", reflect.TypeOf((*MockRaftConfig)(nil).InstanceDirPath))
}

// InstanceId mocks base method.
func (m *MockRaftConfig) InstanceId() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceId")
	ret0, _ := ret[0].(int32)
	return ret0
}

// InstanceId indicates an expected call of InstanceId.
func (mr *MockRaftConfigMockRecorder) InstanceId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceId", reflect.TypeOf((*MockRaftConfig)(nil).InstanceId))
}

// InstanceName mocks base method.
func (m *MockRaftConfig) InstanceName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceName")
	ret0, _ := ret[0].(string)
	return ret0
}

// InstanceName indicates an expected call of InstanceName.
func (mr *MockRaftConfigMockRecorder) InstanceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceName", reflect.TypeOf((*MockRaftConfig)(nil).InstanceName))
}

// LoadConfig mocks base method.
func (m *MockRaftConfig) LoadConfig(arg0 *raft.RaftNode) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadConfig", arg0)
}

// LoadConfig indicates an expected call of LoadConfig.
func (mr *MockRaftConfigMockRecorder) LoadConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockRaftConfig)(nil).LoadConfig), arg0)
}

// LogFilePath mocks base method.
func (m *MockRaftConfig) LogFilePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// LogFilePath indicates an expected call of LogFilePath.
func (mr *MockRaftConfigMockRecorder) LogFilePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogFilePath", reflect.TypeOf((*MockRaftConfig)(nil).LogFilePath))
}

// Logs mocks base method.
func (m *MockRaftConfig) Logs() []*protos.Log {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logs")
	ret0, _ := ret[0].([]*protos.Log)
	return ret0
}

// Logs indicates an expected call of Logs.
func (mr *MockRaftConfigMockRecorder) Logs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logs", reflect.TypeOf((*MockRaftConfig)(nil).Logs))
}

// Peers mocks base method.
func (m *MockRaftConfig) Peers() []raft.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peers")
	ret0, _ := ret[0].([]raft.Peer)
	return ret0
}

// Peers indicates an expected call of Peers.
func (mr *MockRaftConfigMockRecorder) Peers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peers", reflect.TypeOf((*MockRaftConfig)(nil).Peers))
}

// Port mocks base method.
func (m *MockRaftConfig) Port() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Port")
	ret0, _ := ret[0].(string)
	return ret0
}

// Port indicates an expected call of Port.
func (mr *MockRaftConfigMockRecorder) Port() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Port", reflect.TypeOf((*MockRaftConfig)(nil).Port))
}

// SetConfigFilePath mocks base method.
func (m *MockRaftConfig) SetConfigFilePath(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConfigFilePath", arg0)
}

// SetConfigFilePath indicates an expected call of SetConfigFilePath.
func (mr *MockRaftConfigMockRecorder) SetConfigFilePath(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfigFilePath", reflect.TypeOf((*MockRaftConfig)(nil).SetConfigFilePath), arg0)
}

// Version mocks base method.
func (m *MockRaftConfig) Version() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(string)
	return ret0
}

// Version indicates an expected call of Version.
func (mr *MockRaftConfigMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockRaftConfig)(nil).Version))
}

// VotedFor mocks base method.
func (m *MockRaftConfig) VotedFor() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VotedFor")
	ret0, _ := ret[0].(int32)
	return ret0
}

// VotedFor indicates an expected call of VotedFor.
func (mr *MockRaftConfigMockRecorder) VotedFor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VotedFor", reflect.TypeOf((*MockRaftConfig)(nil).VotedFor))
}
