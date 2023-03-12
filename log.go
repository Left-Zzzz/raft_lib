package raftlib

import "time"

// 定义日志类型
type LogType uint8

const (
	// 心跳
	HeartBeat LogType = iota
	// 命令
	LogCommand
	// 空日志补丁，用来覆盖提交未提交的日志，起保护作用
	LogNoOp
)

type Log struct {
	// 记录日志索引号
	Index uint64
	// 记录日志任期号
	Term uint64
	// 记录日志类型
	LogType LogType
	// 记录日志内容
	Data []byte
	// 记录日志首次记录时间，便于follower判断有多少秒的日志需要追踪
	// 为了提高效率，follower日志提交进度追上leader才能处理rpc请求（后期开发）
	AppendedAt time.Time
}
