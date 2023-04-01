package raftlib

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
)

// 暂时用数组缓存，不考虑持久化
type LogStorage struct {
	entries     [][]byte
	commitIndex uint32
	entriesLock sync.RWMutex
	commitLock  sync.Mutex
	// 该map用于存放应用日志的回调函数，方便后期使用该库注册
	callBackFuncMap sync.Map
}

// 返回初始化参数的LogStorage结构体指针
func createLogStorage() *LogStorage {
	return &LogStorage{
		commitIndex: MAX_LOG_INDEX_NUM,
	}
}

func (ls *LogStorage) getCommitIndex() uint32 {
	return atomic.LoadUint32(&ls.commitIndex)
}

// 设置最后一次提交的索引号
func (ls *LogStorage) setCommitIndex(index uint32) {
	atomic.StoreUint32(&ls.commitIndex, index)
}

// 通过日志实体添加日志
func (ls *LogStorage) appendEntryEntity(index uint32, entry Log) error {
	entrEncode, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return ls.appendEntryEncoded(index, entrEncode)
}

// 通过[]byte添加日志
func (ls *LogStorage) appendEntryEncoded(index uint32, entryEncode []byte) error {
	ls.entriesLock.Lock()
	ls.entries = append(ls.entries[:index], entryEncode)
	ls.entriesLock.Unlock()
	entry, err := ls.getEntry(index)
	if err != nil {
		return errors.New("get entry failed")
	}
	logDebug("entry: %v, data:%v", entry, string(entry.Data))

	return nil
}

// 获得日志
func (ls *LogStorage) getEntry(index uint32) (Log, error) {
	ls.entriesLock.RLock()
	defer ls.entriesLock.RUnlock()
	entriesLen := uint32(len(ls.entries))
	logDebug("getEntry(): entriesLen:%v, index:%v", entriesLen, index)
	if entriesLen <= index {
		return Log{}, errors.New("param \"index\" is larger than len of entries")
	}
	rowData := ls.entries[index]

	logEntry := Log{}
	json.Unmarshal(rowData, &logEntry)
	return logEntry, nil
}

// 判断日志项索引号和任期是否正确
func (ls *LogStorage) isIdxTermCorrect(index uint32, term uint64) bool {
	// 如果为没有上一个日志项，则不用检验任期，返回正确
	if index == MAX_LOG_INDEX_NUM {
		logDebug("isIdxTermCorrect(): prevLogIndex is MAX_LOG_INDEX_NUM, return.")
		return true
	}
	entry, err := ls.getEntry(index)
	if err != nil {
		logError("an error occurred while getting entry:%v, index: %v", err, index)
		return false
	}
	return entry.Term == term
}

// 提交日志，调用回调函数处理日志命令并执行
func (ls *LogStorage) commit(index uint32, execFunc func([]byte) error) error {
	// 获取日志
	logEntry, err := ls.getEntry(index)
	// 如果获取日志出错，返回错误
	if err != nil {
		return err
	}
	// 如果不是LogCommand类型，直接跳过
	if logEntry.LogType != LogCommand {
		return nil
	}
	// 调用回调函数
	err = execFunc(logEntry.Data)
	// 如果不发生错误，更新提交日志索引
	if err == nil {
		ls.setCommitIndex(index)
		logDebug("LogStorage.commit(): commitIndex: %v", index)
	}
	return err
}

// 提交日志，调用回调函数处理日志命令并执行
func (ls *LogStorage) batchCommit(index uint32, execFunc func([]byte) error) error {
	// 如果没有日志提交，直接返回
	if index == MAX_LOG_INDEX_NUM {
		return nil
	}
	ok := ls.commitLock.TryLock()
	// 如果已经被获取，停止执行，每时刻只能有一个执行体执行提交日志操作
	if !ok {
		return nil
	}
	defer ls.commitLock.Unlock()
	for curIndex := genNextLogIndex(ls.getCommitIndex()); curIndex <= index; curIndex++ {
		logDebug("commit entry(index:%v), commitIndex:%v", curIndex, index)
		err := ls.commit(curIndex, execFunc)
		// 如果获取日志出错，返回错误
		if err != nil {
			return err
		}
	}
	return nil
}

// 注册应用日志的回调函数
func (ls *LogStorage) registerCallBackFunc(f func([]byte) error) {
	ls.callBackFuncMap.Store(execCommandFuncName, f)
}
