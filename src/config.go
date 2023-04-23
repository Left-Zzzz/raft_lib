package raftlib

import (
	"time"
)

type Config struct {
	// 是否启用日志
	DEBUG bool

	// 设置超时时间
	// 心跳超时时间
	HeartbeatTimeout time.Duration
	// 选举超时时间
	ElectionTimeout time.Duration
	// RPC响应超时时间
	RpcTimeout time.Duration

	// Servers: 服务器元信息列表
	Servers []Server

	// 当前Raft节点
	LocalServer Server
}

func CreateConfig() *Config {
	config := &Config{
		DEBUG:            true,
		HeartbeatTimeout: 300 * time.Millisecond,
		ElectionTimeout:  200 * time.Millisecond,
		RpcTimeout:       100 * time.Millisecond,
		Servers: []Server{
			{Voter, "0", "127.0.0.1", "5676"},
			{Voter, "1", "127.0.0.1", "5677"},
		},
		LocalServer: Server{Voter, "0", "127.0.0.1", "5676"},
	}
	return config
}

const (
	// rpc proto版本号
	RPOTO_VER_REQUEST_VOTE_REQUEST  uint32 = 1
	RPOTO_VER_REQUEST_VOTE_RESPONSE uint32 = 1
	RPOTO_VER_APPEND_ENTRY_REQUEST  uint32 = 1
	RPOTO_VER_APPEND_ENTRY_RESPONSE uint32 = 1
	PROTO_VER_EXEC_COMMAND_REQUEST  uint32 = 1
	PROTO_VER_EXEC_COMMAND_RESPONSE uint32 = 1

	// 日志相关
	// 设置最大日志项数量
	MAX_LOG_INDEX_NUM = 0xFFFFFFFF
	// 定义应用日志的回调函数名字
	EXEC_COMMAND_FUNC_NAME string = "execCommandFunc"

	// ServerSuffrage类型(暂不使用)
	// Voter：集群成员变更时，当节点日志纪录追上领导者节点时候，具有Voter权限
	Voter ServerSuffrage = 0
	// Nonvoter：集群成员变更时，当节点正在纪录日志且日志未追上领导者节点时，不具有Voter权限
	Nonvoter ServerSuffrage = 1
)

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
