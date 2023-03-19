package raftlib

import (
	"encoding/json"
	"sync"
)

// 暂时用数组缓存，不考虑持久化
type LogStorage struct {
	entries     [][]byte
	entriesLock sync.RWMutex
	// 该map用于存放应用日志的回调函数，方便后期使用该库注册
	callBackFuncMap sync.Map
}

// 添加日志
func (ls *LogStorage) appendEntry(entry Log) error {
	entry_encode, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	ls.entriesLock.Lock()
	defer ls.entriesLock.Unlock()
	ls.entries = append(ls.entries, entry_encode)
	logDebug("l.entries: %v", ls.entries)
	return nil
}

// 获得日志
func (ls *LogStorage) getEntry(index uint32) Log {
	ls.entriesLock.RLock()
	rowData := ls.entries[index]
	ls.entriesLock.RUnlock()
	logEntry := Log{}
	json.Unmarshal(rowData, &logEntry)
	return logEntry
}

// 提交日志，调用回调函数处理日志命令并执行
func (ls *LogStorage) commit(index uint32, execfunc func([]byte) error) error {
	logEntry := &Log{}
	ls.entriesLock.RLock()
	err := json.Unmarshal(ls.entries[index], logEntry)
	ls.entriesLock.RUnlock()
	if err != nil {
		logWarn("an error happened while decoding logEntry from logStorage: %v", err)
	}
	if logEntry.LogType != LogCommand {
		return nil
	}
	return execfunc(logEntry.Data)
}

// 注册应用日志的回调函数
func (ls *LogStorage) registerCallBackFunc(f func([]byte) error) {
	ls.callBackFuncMap.Store(execCommandFuncName, f)
}
