package raftlib

import "time"

const (
	// 心跳超时时间
	HeartbeatTimeout = 1000 * time.Millisecond
	// 选举超时时间
	ElectionTimeout = 1000 * time.Millisecond
)

// Servers: 服务器元信息列表
var Servers = []Server{
	{1, "0", "192.168.100.3"},
	{1, "1", "192.168.100.4"},
}

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
	// Voter：当节点日志纪录追上领导者节点时候，具有Voter权限
	Voter ServerSuffrage = iota
	// Nonvoter：当节点正在纪录日志且日志未追上领导者节点时，不具有Voter权限
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
}
