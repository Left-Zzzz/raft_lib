package raftlib

import (
	"time"
)

// 设置超时时间
const (
	// 心跳超时时间
	HeartbeatTimeout = 1000 * time.Millisecond
	// 选举超时时间
	ElectionTimeout = 1000 * time.Millisecond
	// RPC响应超时时间
	RpcTimeout = 50000000 * time.Millisecond
)

// 设置最大数量
const (
	MAX_LOG_INDEX_NUM = 0xFFFFFFFF
)

// Servers: 服务器元信息列表
var Servers = []Server{
	{Voter, "0", "127.0.0.1", "5676"},
	{Voter, "1", "127.0.0.1", "5677"},
}

// 记录本地服务器
var localServer Server = Servers[0]

// ClientAddress : 网络ip地址
type ClientAddress string

// ClientPort : 网络端口
type ClientPort string

// ServerAddress : 网络ip地址
type ServerAddress string

// ServerPort : 网络端口
type ServerPort string

// ServerID : 服务器id，总是唯一确定的
type ServerID string

// ServerSuffrage: 表示节点是否获得投票
type ServerSuffrage int

// ServerSuffrage类型
const (
	// Voter：集群成员变更时，当节点日志纪录追上领导者节点时候，具有Voter权限
	Voter ServerSuffrage = iota
	// Nonvoter：集群成员变更时，当节点正在纪录日志且日志未追上领导者节点时，不具有Voter权限
	Nonvoter
)

// Server: 服务器元信息
type Server struct {
	// Suffrage: 表示节点是否获得投票
	Suffrage ServerSuffrage
	// ID: 总是唯一确定的服务器标识码
	ID ServerID
	// Address: 服务器网络地址
	Address ServerAddress
	// Port: 服务器端口号
	Port ServerPort
}

// rpc proto版本号
var (
	RPOTO_VER_REQUEST_VOTE_REQUEST  uint32 = 1
	RPOTO_VER_REQUEST_VOTE_RESPONSE uint32 = 1
	RPOTO_VER_APPEND_ENTRY_REQUEST  uint32 = 1
	RPOTO_VER_APPEND_ENTRY_RESPONSE uint32 = 1
	PROTO_VER_EXEC_COMMAND_REQUEST  uint32 = 1
	PROTO_VER_EXEC_COMMAND_RESPONSE uint32 = 1
)

// 是否启用Debug日志
const (
	DEBUG = true
)

// Dubug中需要打印的变量
var (
	isPrintRpcAppendEntryRequestCh = [2]bool{false, false}
)

// 定义应用日志的回调函数名字
const (
	execCommandFuncName string = "execCommandFunc"
)
