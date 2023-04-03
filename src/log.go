package raftlib

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
	// 日志索引号，也可通过entries数组的下标查找索引号
	Index uint32
	// 记录日志任期号
	Term uint64
	// 记录日志类型
	LogType LogType
	// 记录日志内容
	Data []byte
	// 记录日志首次记录时间，便于follower判断有多少秒的日志需要追踪
	// 为了提高效率，follower日志提交进度追上leader才能处理rpc请求（后期开发）
	// AppendedAt time.Time
}

func createLogEntry(index uint32, logType LogType, term uint64, data []byte) Log {
	return Log{
		Index:   index,
		Term:    term,
		LogType: logType,
		Data:    data,
	}
}
