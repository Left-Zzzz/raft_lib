// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.0
// source: messageDefine.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 选举请求RPC
type RequestVoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver         *uint32 `protobuf:"varint,1,opt,name=Ver,proto3,oneof" json:"Ver,omitempty"`
	Term        uint64  `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	CandidateId string  `protobuf:"bytes,3,opt,name=CandidateId,proto3" json:"CandidateId,omitempty"`
	LastLogIdx  uint32  `protobuf:"varint,4,opt,name=LastLogIdx,proto3" json:"LastLogIdx,omitempty"`
	LastLogTerm uint64  `protobuf:"varint,5,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
}

func (x *RequestVoteRequest) Reset() {
	*x = RequestVoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messageDefine_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteRequest) ProtoMessage() {}

func (x *RequestVoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messageDefine_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteRequest.ProtoReflect.Descriptor instead.
func (*RequestVoteRequest) Descriptor() ([]byte, []int) {
	return file_messageDefine_proto_rawDescGZIP(), []int{0}
}

func (x *RequestVoteRequest) GetVer() uint32 {
	if x != nil && x.Ver != nil {
		return *x.Ver
	}
	return 0
}

func (x *RequestVoteRequest) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteRequest) GetCandidateId() string {
	if x != nil {
		return x.CandidateId
	}
	return ""
}

func (x *RequestVoteRequest) GetLastLogIdx() uint32 {
	if x != nil {
		return x.LastLogIdx
	}
	return 0
}

func (x *RequestVoteRequest) GetLastLogTerm() uint64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

// 选举请求回复RPC
type RequestVoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver     *uint32 `protobuf:"varint,1,opt,name=Ver,proto3,oneof" json:"Ver,omitempty"`
	VoterID string  `protobuf:"bytes,2,opt,name=VoterID,proto3" json:"VoterID,omitempty"`
	Term    uint64  `protobuf:"varint,3,opt,name=Term,proto3" json:"Term,omitempty"`
	// 安全补丁，保证leader日志时最新的
	// candidate 收到比他日志还新的回复就会将自己转换为follower
	LastLogIdx  uint32 `protobuf:"varint,4,opt,name=LastLogIdx,proto3" json:"LastLogIdx,omitempty"`
	LastLogTerm uint64 `protobuf:"varint,5,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
	VoteGranted bool   `protobuf:"varint,6,opt,name=VoteGranted,proto3" json:"VoteGranted,omitempty"`
}

func (x *RequestVoteResponse) Reset() {
	*x = RequestVoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messageDefine_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteResponse) ProtoMessage() {}

func (x *RequestVoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messageDefine_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteResponse.ProtoReflect.Descriptor instead.
func (*RequestVoteResponse) Descriptor() ([]byte, []int) {
	return file_messageDefine_proto_rawDescGZIP(), []int{1}
}

func (x *RequestVoteResponse) GetVer() uint32 {
	if x != nil && x.Ver != nil {
		return *x.Ver
	}
	return 0
}

func (x *RequestVoteResponse) GetVoterID() string {
	if x != nil {
		return x.VoterID
	}
	return ""
}

func (x *RequestVoteResponse) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteResponse) GetLastLogIdx() uint32 {
	if x != nil {
		return x.LastLogIdx
	}
	return 0
}

func (x *RequestVoteResponse) GetLastLogTerm() uint64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

func (x *RequestVoteResponse) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

// appendEntryRequest 参数意义参考Raft论文
type AppendEntryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver          *uint32 `protobuf:"varint,1,opt,name=Ver,proto3,oneof" json:"Ver,omitempty"`
	LeaderTerm   uint64  `protobuf:"varint,2,opt,name=LeaderTerm,proto3" json:"LeaderTerm,omitempty"`
	LeaderID     string  `protobuf:"bytes,3,opt,name=LeaderID,proto3" json:"LeaderID,omitempty"`
	PrevLogIndex uint32  `protobuf:"varint,4,opt,name=PrevLogIndex,proto3" json:"PrevLogIndex,omitempty"`
	PrevLogTerm  uint64  `protobuf:"varint,5,opt,name=PrevLogTerm,proto3" json:"PrevLogTerm,omitempty"`
	LogType      uint32  `protobuf:"varint,6,opt,name=LogType,proto3" json:"LogType,omitempty"`
	Entry        []byte  `protobuf:"bytes,7,opt,name=Entry,proto3,oneof" json:"Entry,omitempty"`
	LeaderCommit uint32  `protobuf:"varint,8,opt,name=LeaderCommit,proto3" json:"LeaderCommit,omitempty"`
}

func (x *AppendEntryRequest) Reset() {
	*x = AppendEntryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messageDefine_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntryRequest) ProtoMessage() {}

func (x *AppendEntryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messageDefine_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntryRequest.ProtoReflect.Descriptor instead.
func (*AppendEntryRequest) Descriptor() ([]byte, []int) {
	return file_messageDefine_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntryRequest) GetVer() uint32 {
	if x != nil && x.Ver != nil {
		return *x.Ver
	}
	return 0
}

func (x *AppendEntryRequest) GetLeaderTerm() uint64 {
	if x != nil {
		return x.LeaderTerm
	}
	return 0
}

func (x *AppendEntryRequest) GetLeaderID() string {
	if x != nil {
		return x.LeaderID
	}
	return ""
}

func (x *AppendEntryRequest) GetPrevLogIndex() uint32 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendEntryRequest) GetPrevLogTerm() uint64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendEntryRequest) GetLogType() uint32 {
	if x != nil {
		return x.LogType
	}
	return 0
}

func (x *AppendEntryRequest) GetEntry() []byte {
	if x != nil {
		return x.Entry
	}
	return nil
}

func (x *AppendEntryRequest) GetLeaderCommit() uint32 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

// appendEntryResponse 参数意义参考Raft论文
type AppendEntryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver     *uint32 `protobuf:"varint,1,opt,name=Ver,proto3,oneof" json:"Ver,omitempty"`
	Term    uint64  `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	Success bool    `protobuf:"varint,3,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *AppendEntryResponse) Reset() {
	*x = AppendEntryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messageDefine_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntryResponse) ProtoMessage() {}

func (x *AppendEntryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messageDefine_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntryResponse.ProtoReflect.Descriptor instead.
func (*AppendEntryResponse) Descriptor() ([]byte, []int) {
	return file_messageDefine_proto_rawDescGZIP(), []int{3}
}

func (x *AppendEntryResponse) GetVer() uint32 {
	if x != nil && x.Ver != nil {
		return *x.Ver
	}
	return 0
}

func (x *AppendEntryResponse) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntryResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// 客户端向Raft集群发送执行命令请求，请求报文
type ExecCommandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver     *uint32 `protobuf:"varint,1,opt,name=Ver,proto3,oneof" json:"Ver,omitempty"`
	LogType uint32  `protobuf:"varint,2,opt,name=LogType,proto3" json:"LogType,omitempty"`
	Command []byte  `protobuf:"bytes,3,opt,name=Command,proto3" json:"Command,omitempty"`
}

func (x *ExecCommandRequest) Reset() {
	*x = ExecCommandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messageDefine_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecCommandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecCommandRequest) ProtoMessage() {}

func (x *ExecCommandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messageDefine_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecCommandRequest.ProtoReflect.Descriptor instead.
func (*ExecCommandRequest) Descriptor() ([]byte, []int) {
	return file_messageDefine_proto_rawDescGZIP(), []int{4}
}

func (x *ExecCommandRequest) GetVer() uint32 {
	if x != nil && x.Ver != nil {
		return *x.Ver
	}
	return 0
}

func (x *ExecCommandRequest) GetLogType() uint32 {
	if x != nil {
		return x.LogType
	}
	return 0
}

func (x *ExecCommandRequest) GetCommand() []byte {
	if x != nil {
		return x.Command
	}
	return nil
}

// 客户端向Raft集群发送执行命令请求，响应报文
type ExecCommandResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver           *uint32 `protobuf:"varint,1,opt,name=Ver,proto3,oneof" json:"Ver,omitempty"`
	LeaderID      string  `protobuf:"bytes,2,opt,name=LeaderID,proto3" json:"LeaderID,omitempty"`
	LeaderAddress string  `protobuf:"bytes,3,opt,name=LeaderAddress,proto3" json:"LeaderAddress,omitempty"`
	LeaderPort    string  `protobuf:"bytes,4,opt,name=LeaderPort,proto3" json:"LeaderPort,omitempty"`
	Success       bool    `protobuf:"varint,5,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *ExecCommandResponse) Reset() {
	*x = ExecCommandResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messageDefine_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecCommandResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecCommandResponse) ProtoMessage() {}

func (x *ExecCommandResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messageDefine_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecCommandResponse.ProtoReflect.Descriptor instead.
func (*ExecCommandResponse) Descriptor() ([]byte, []int) {
	return file_messageDefine_proto_rawDescGZIP(), []int{5}
}

func (x *ExecCommandResponse) GetVer() uint32 {
	if x != nil && x.Ver != nil {
		return *x.Ver
	}
	return 0
}

func (x *ExecCommandResponse) GetLeaderID() string {
	if x != nil {
		return x.LeaderID
	}
	return ""
}

func (x *ExecCommandResponse) GetLeaderAddress() string {
	if x != nil {
		return x.LeaderAddress
	}
	return ""
}

func (x *ExecCommandResponse) GetLeaderPort() string {
	if x != nil {
		return x.LeaderPort
	}
	return ""
}

func (x *ExecCommandResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_messageDefine_proto protoreflect.FileDescriptor

var file_messageDefine_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x69, 0x62, 0x22, 0xab,
	0x01, 0x0a, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x03, 0x56, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x48, 0x00, 0x52, 0x03, 0x56, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x12, 0x0a, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d,
	0x12, 0x20, 0x0a, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x64, 0x78,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49,
	0x64, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67,
	0x54, 0x65, 0x72, 0x6d, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x56, 0x65, 0x72, 0x22, 0xc6, 0x01, 0x0a,
	0x13, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x56, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x48, 0x00, 0x52, 0x03, 0x56, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07, 0x56,
	0x6f, 0x74, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x56, 0x6f,
	0x74, 0x65, 0x72, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x4c, 0x61, 0x73,
	0x74, 0x4c, 0x6f, 0x67, 0x49, 0x64, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x4c,
	0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x64, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61, 0x73,
	0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x56,
	0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x42, 0x06, 0x0a,
	0x04, 0x5f, 0x56, 0x65, 0x72, 0x22, 0x98, 0x02, 0x0a, 0x12, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x03,
	0x56, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x03, 0x56, 0x65, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x0a, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x54,
	0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12,
	0x22, 0x0a, 0x0c, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65,
	0x72, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f,
	0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x4c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x19, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x01,
	0x52, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0c, 0x4c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0c, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x42, 0x06,
	0x0a, 0x04, 0x5f, 0x56, 0x65, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x22, 0x62, 0x0a, 0x13, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x56, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x03, 0x56, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x12,
	0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65,
	0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x06, 0x0a, 0x04,
	0x5f, 0x56, 0x65, 0x72, 0x22, 0x67, 0x0a, 0x12, 0x65, 0x78, 0x65, 0x63, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x03, 0x56, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x03, 0x56, 0x65, 0x72, 0x88, 0x01,
	0x01, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x4c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x56, 0x65, 0x72, 0x22, 0xb0, 0x01,
	0x0a, 0x13, 0x65, 0x78, 0x65, 0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x56, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x48, 0x00, 0x52, 0x03, 0x56, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1a, 0x0a, 0x08,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x4c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1e,
	0x0a, 0x0a, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x56, 0x65, 0x72,
	0x32, 0xea, 0x01, 0x0a, 0x0a, 0x52, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x48, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x1b,
	0x2e, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x69, 0x62, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x6c, 0x69, 0x62, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0b, 0x41, 0x70, 0x70,
	0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x6c,
	0x69, 0x62, 0x2e, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x69, 0x62, 0x2e,
	0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0b, 0x45, 0x78, 0x65, 0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x12, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x69, 0x62, 0x2e, 0x65, 0x78, 0x65,
	0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x6c, 0x69, 0x62, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x06, 0x5a,
	0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messageDefine_proto_rawDescOnce sync.Once
	file_messageDefine_proto_rawDescData = file_messageDefine_proto_rawDesc
)

func file_messageDefine_proto_rawDescGZIP() []byte {
	file_messageDefine_proto_rawDescOnce.Do(func() {
		file_messageDefine_proto_rawDescData = protoimpl.X.CompressGZIP(file_messageDefine_proto_rawDescData)
	})
	return file_messageDefine_proto_rawDescData
}

var file_messageDefine_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_messageDefine_proto_goTypes = []interface{}{
	(*RequestVoteRequest)(nil),  // 0: raftlib.requestVoteRequest
	(*RequestVoteResponse)(nil), // 1: raftlib.requestVoteResponse
	(*AppendEntryRequest)(nil),  // 2: raftlib.appendEntryRequest
	(*AppendEntryResponse)(nil), // 3: raftlib.appendEntryResponse
	(*ExecCommandRequest)(nil),  // 4: raftlib.execCommandRequest
	(*ExecCommandResponse)(nil), // 5: raftlib.execCommandResponse
}
var file_messageDefine_proto_depIdxs = []int32{
	0, // 0: raftlib.RpcService.RequestVote:input_type -> raftlib.requestVoteRequest
	2, // 1: raftlib.RpcService.AppendEntry:input_type -> raftlib.appendEntryRequest
	4, // 2: raftlib.RpcService.ExecCommand:input_type -> raftlib.execCommandRequest
	1, // 3: raftlib.RpcService.RequestVote:output_type -> raftlib.requestVoteResponse
	3, // 4: raftlib.RpcService.AppendEntry:output_type -> raftlib.appendEntryResponse
	5, // 5: raftlib.RpcService.ExecCommand:output_type -> raftlib.execCommandResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_messageDefine_proto_init() }
func file_messageDefine_proto_init() {
	if File_messageDefine_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messageDefine_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteRequest); i {
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
		file_messageDefine_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteResponse); i {
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
		file_messageDefine_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntryRequest); i {
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
		file_messageDefine_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntryResponse); i {
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
		file_messageDefine_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecCommandRequest); i {
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
		file_messageDefine_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecCommandResponse); i {
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
	file_messageDefine_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_messageDefine_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_messageDefine_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_messageDefine_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_messageDefine_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_messageDefine_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messageDefine_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_messageDefine_proto_goTypes,
		DependencyIndexes: file_messageDefine_proto_depIdxs,
		MessageInfos:      file_messageDefine_proto_msgTypes,
	}.Build()
	File_messageDefine_proto = out.File
	file_messageDefine_proto_rawDesc = nil
	file_messageDefine_proto_goTypes = nil
	file_messageDefine_proto_depIdxs = nil
}
