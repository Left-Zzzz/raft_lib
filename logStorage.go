package raftlib

import (
	"encoding/json"
	"sync"
)

// 暂时用数组缓存，不考虑持久化
type LogStorage struct {
	entries     [][]byte
	entriesLock sync.RWMutex
}

// 添加日志
func (l *LogStorage) appendEntry(entry Log) error {
	entry_encode, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	l.entriesLock.Lock()
	defer l.entriesLock.Unlock()
	l.entries = append(l.entries, entry_encode)
	logDebug("l.entries: %v", l.entries)
	return nil
}

// 获得日志
func (l *LogStorage) getEntry(index uint32) Log {
	l.entriesLock.RLock()
	rowData := l.entries[index]
	l.entriesLock.RUnlock()
	logEntry := Log{}
	json.Unmarshal(rowData, &logEntry)
	return logEntry
}

// 提交日志，调用回调函数处理日志命令并执行
func (l *LogStorage) commit(index uint32, execfunc func([]byte) error) error {
	logEntry := &Log{}
	l.entriesLock.RLock()
	err := json.Unmarshal(l.entries[index], logEntry)
	l.entriesLock.RUnlock()
	if err != nil {
		logWarn("an error happened while decoding logEntry from logStorage: %v", err)
	}
	if logEntry.LogType != LogCommand {
		return nil
	}
	return execfunc(logEntry.Data)
}
