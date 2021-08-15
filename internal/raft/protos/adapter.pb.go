// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: internal/raft/protos/adapter.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId      int32 `protobuf:"varint,1,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	CurrentTerm int32 `protobuf:"varint,2,opt,name=CurrentTerm,proto3" json:"CurrentTerm,omitempty"`
	LogLength   int32 `protobuf:"varint,3,opt,name=LogLength,proto3" json:"LogLength,omitempty"`
	LastTerm    int32 `protobuf:"varint,4,opt,name=LastTerm,proto3" json:"LastTerm,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{0}
}

func (x *VoteRequest) GetNodeId() int32 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *VoteRequest) GetCurrentTerm() int32 {
	if x != nil {
		return x.CurrentTerm
	}
	return 0
}

func (x *VoteRequest) GetLogLength() int32 {
	if x != nil {
		return x.LogLength
	}
	return 0
}

func (x *VoteRequest) GetLastTerm() int32 {
	if x != nil {
		return x.LastTerm
	}
	return 0
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId      int32 `protobuf:"varint,1,opt,name=PeerId,proto3" json:"PeerId,omitempty"`
	Term        int32 `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	VoteGranted bool  `protobuf:"varint,3,opt,name=VoteGranted,proto3" json:"VoteGranted,omitempty"`
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{1}
}

func (x *VoteResponse) GetPeerId() int32 {
	if x != nil {
		return x.PeerId
	}
	return 0
}

func (x *VoteResponse) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteResponse) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

type LogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeaderId     int32  `protobuf:"varint,1,opt,name=LeaderId,proto3" json:"LeaderId,omitempty"`
	CurrentTerm  int32  `protobuf:"varint,2,opt,name=CurrentTerm,proto3" json:"CurrentTerm,omitempty"`
	SentLength   int32  `protobuf:"varint,3,opt,name=SentLength,proto3" json:"SentLength,omitempty"`
	PrevLogTerm  int32  `protobuf:"varint,4,opt,name=PrevLogTerm,proto3" json:"PrevLogTerm,omitempty"`
	CommitLength int32  `protobuf:"varint,5,opt,name=CommitLength,proto3" json:"CommitLength,omitempty"`
	Entries      []*Log `protobuf:"bytes,6,rep,name=Entries,proto3" json:"Entries,omitempty"`
}

func (x *LogRequest) Reset() {
	*x = LogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogRequest) ProtoMessage() {}

func (x *LogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogRequest.ProtoReflect.Descriptor instead.
func (*LogRequest) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{2}
}

func (x *LogRequest) GetLeaderId() int32 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *LogRequest) GetCurrentTerm() int32 {
	if x != nil {
		return x.CurrentTerm
	}
	return 0
}

func (x *LogRequest) GetSentLength() int32 {
	if x != nil {
		return x.SentLength
	}
	return 0
}

func (x *LogRequest) GetPrevLogTerm() int32 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *LogRequest) GetCommitLength() int32 {
	if x != nil {
		return x.CommitLength
	}
	return 0
}

func (x *LogRequest) GetEntries() []*Log {
	if x != nil {
		return x.Entries
	}
	return nil
}

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term    int32            `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	Message *structpb.Struct `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{3}
}

func (x *Log) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Log) GetMessage() *structpb.Struct {
	if x != nil {
		return x.Message
	}
	return nil
}

type LogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FollowerId int32 `protobuf:"varint,1,opt,name=FollowerId,proto3" json:"FollowerId,omitempty"`
	Term       int32 `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	AckLength  int32 `protobuf:"varint,3,opt,name=AckLength,proto3" json:"AckLength,omitempty"`
	Success    bool  `protobuf:"varint,4,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *LogResponse) Reset() {
	*x = LogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogResponse) ProtoMessage() {}

func (x *LogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogResponse.ProtoReflect.Descriptor instead.
func (*LogResponse) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{4}
}

func (x *LogResponse) GetFollowerId() int32 {
	if x != nil {
		return x.FollowerId
	}
	return 0
}

func (x *LogResponse) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogResponse) GetAckLength() int32 {
	if x != nil {
		return x.AckLength
	}
	return 0
}

func (x *LogResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type BroadcastMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *BroadcastMessageResponse) Reset() {
	*x = BroadcastMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastMessageResponse) ProtoMessage() {}

func (x *BroadcastMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastMessageResponse.ProtoReflect.Descriptor instead.
func (*BroadcastMessageResponse) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{5}
}

func (x *BroadcastMessageResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type RaftPersistentState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RaftLogs         []*Log `protobuf:"bytes,1,rep,name=RaftLogs,proto3" json:"RaftLogs,omitempty"`
	RaftCurrentTerm  int32  `protobuf:"varint,2,opt,name=RaftCurrentTerm,proto3" json:"RaftCurrentTerm,omitempty"`
	RaftVotedFor     int32  `protobuf:"varint,3,opt,name=RaftVotedFor,proto3" json:"RaftVotedFor,omitempty"`
	RaftCommitLength int32  `protobuf:"varint,4,opt,name=RaftCommitLength,proto3" json:"RaftCommitLength,omitempty"`
}

func (x *RaftPersistentState) Reset() {
	*x = RaftPersistentState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_raft_protos_adapter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftPersistentState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftPersistentState) ProtoMessage() {}

func (x *RaftPersistentState) ProtoReflect() protoreflect.Message {
	mi := &file_internal_raft_protos_adapter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftPersistentState.ProtoReflect.Descriptor instead.
func (*RaftPersistentState) Descriptor() ([]byte, []int) {
	return file_internal_raft_protos_adapter_proto_rawDescGZIP(), []int{6}
}

func (x *RaftPersistentState) GetRaftLogs() []*Log {
	if x != nil {
		return x.RaftLogs
	}
	return nil
}

func (x *RaftPersistentState) GetRaftCurrentTerm() int32 {
	if x != nil {
		return x.RaftCurrentTerm
	}
	return 0
}

func (x *RaftPersistentState) GetRaftVotedFor() int32 {
	if x != nil {
		return x.RaftVotedFor
	}
	return 0
}

func (x *RaftPersistentState) GetRaftCommitLength() int32 {
	if x != nil {
		return x.RaftCommitLength
	}
	return 0
}

var File_internal_raft_protos_adapter_proto protoreflect.FileDescriptor

var file_internal_raft_protos_adapter_proto_rawDesc = []byte{
	0x0a, 0x22, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x64, 0x61, 0x70, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x1c, 0x0a, 0x09,
	0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x61,
	0x73, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x4c, 0x61,
	0x73, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x5c, 0x0a, 0x0c, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x54, 0x65,
	0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61,
	0x6e, 0x74, 0x65, 0x64, 0x22, 0xd0, 0x01, 0x0a, 0x0a, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x20, 0x0a, 0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x72,
	0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x53, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54,
	0x65, 0x72, 0x6d, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x4c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x1e, 0x0a, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x07,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x4c, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12, 0x12,
	0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x54, 0x65,
	0x72, 0x6d, 0x12, 0x31, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x79, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x63, 0x6b, 0x4c,
	0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x41, 0x63, 0x6b,
	0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x34, 0x0a, 0x18, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xb1, 0x01, 0x0a, 0x13, 0x52, 0x61, 0x66, 0x74, 0x50,
	0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20,
	0x0a, 0x08, 0x52, 0x61, 0x66, 0x74, 0x4c, 0x6f, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x04, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x08, 0x52, 0x61, 0x66, 0x74, 0x4c, 0x6f, 0x67, 0x73,
	0x12, 0x28, 0x0a, 0x0f, 0x52, 0x61, 0x66, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54,
	0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x52, 0x61, 0x66, 0x74, 0x43,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x22, 0x0a, 0x0c, 0x52, 0x61,
	0x66, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x64, 0x46, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x52, 0x61, 0x66, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x64, 0x46, 0x6f, 0x72, 0x12, 0x2a,
	0x0a, 0x10, 0x52, 0x61, 0x66, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x4c, 0x65, 0x6e, 0x67,
	0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x52, 0x61, 0x66, 0x74, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x32, 0xad, 0x01, 0x0a, 0x04, 0x52,
	0x61, 0x66, 0x74, 0x12, 0x32, 0x0a, 0x13, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f,
	0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x50, 0x65, 0x65, 0x72, 0x12, 0x0c, 0x2e, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x12, 0x0b, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x46, 0x0a, 0x10, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x1a,
	0x19, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x61, 0x72, 0x61, 0x70, 0x68, 0x69,
	0x78, 0x2f, 0x53, 0x69, 0x66, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x72,
	0x61, 0x66, 0x74, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x61, 0x64, 0x61, 0x70, 0x74, 0x65, 0x72, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_raft_protos_adapter_proto_rawDescOnce sync.Once
	file_internal_raft_protos_adapter_proto_rawDescData = file_internal_raft_protos_adapter_proto_rawDesc
)

func file_internal_raft_protos_adapter_proto_rawDescGZIP() []byte {
	file_internal_raft_protos_adapter_proto_rawDescOnce.Do(func() {
		file_internal_raft_protos_adapter_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_raft_protos_adapter_proto_rawDescData)
	})
	return file_internal_raft_protos_adapter_proto_rawDescData
}

var file_internal_raft_protos_adapter_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_raft_protos_adapter_proto_goTypes = []interface{}{
	(*VoteRequest)(nil),              // 0: VoteRequest
	(*VoteResponse)(nil),             // 1: VoteResponse
	(*LogRequest)(nil),               // 2: LogRequest
	(*Log)(nil),                      // 3: Log
	(*LogResponse)(nil),              // 4: LogResponse
	(*BroadcastMessageResponse)(nil), // 5: BroadcastMessageResponse
	(*RaftPersistentState)(nil),      // 6: RaftPersistentState
	(*structpb.Struct)(nil),          // 7: google.protobuf.Struct
}
var file_internal_raft_protos_adapter_proto_depIdxs = []int32{
	3, // 0: LogRequest.Entries:type_name -> Log
	7, // 1: Log.Message:type_name -> google.protobuf.Struct
	3, // 2: RaftPersistentState.RaftLogs:type_name -> Log
	0, // 3: Raft.RequestVoteFromPeer:input_type -> VoteRequest
	2, // 4: Raft.ReplicateLog:input_type -> LogRequest
	7, // 5: Raft.BroadcastMessage:input_type -> google.protobuf.Struct
	1, // 6: Raft.RequestVoteFromPeer:output_type -> VoteResponse
	4, // 7: Raft.ReplicateLog:output_type -> LogResponse
	5, // 8: Raft.BroadcastMessage:output_type -> BroadcastMessageResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_raft_protos_adapter_proto_init() }
func file_internal_raft_protos_adapter_proto_init() {
	if File_internal_raft_protos_adapter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_raft_protos_adapter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_raft_protos_adapter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_raft_protos_adapter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_raft_protos_adapter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_raft_protos_adapter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_raft_protos_adapter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastMessageResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_raft_protos_adapter_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RaftPersistentState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_raft_protos_adapter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_raft_protos_adapter_proto_goTypes,
		DependencyIndexes: file_internal_raft_protos_adapter_proto_depIdxs,
		MessageInfos:      file_internal_raft_protos_adapter_proto_msgTypes,
	}.Build()
	File_internal_raft_protos_adapter_proto = out.File
	file_internal_raft_protos_adapter_proto_rawDesc = nil
	file_internal_raft_protos_adapter_proto_goTypes = nil
	file_internal_raft_protos_adapter_proto_depIdxs = nil
}
