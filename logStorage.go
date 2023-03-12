package raftlib

// 暂时用数组缓存，不考虑持久化
type LogStorage struct {
	entries [][]byte
}
