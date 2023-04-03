package raftlib

import (
	"sync"
	"sync/atomic"
)

// 状态机角色：follower、candidate、leader
type RaftState uint32

// 状态机角色常量定义
const (
	// Follower 追随者，状态机默认角色
	Follower RaftState = iota

	// Candidate 候选者
	Candidate

	// Leader 领导者
	Leader
)

// 状态角色信息转为字符串形式
func (s RaftState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unkonwn"
	}
}

// raftState：维护raft节点状态变量以及
// 为节点提供线程安全的set/get方法
// 详见论文——Raft 一致性算法章节
type raftState struct {
	// 当前任期号，初始值0，单增
	currentTerm uint64

	// 最大日志提交索引号，初始值MAX_LOG_INDEX_NUM(-1)，单增
	commitIndex uint32

	// 保护一下两个参数的读写操作
	lastLock sync.Mutex

	// 缓存最新log Index(初始值：MAX_LOG_INDEX_NUM(-1))/Term
	lastLogIndex uint32
	lastLogTerm  uint64

	// 已投票给候选者中候选者最大任期号
	latestVoteGrantedTerm uint64

	// 记录正在运行的goroutines数量
	routinesGroup sync.WaitGroup

	// 当前状态机角色
	state RaftState
}

// 获取当前角色
func (r *raftState) getState() RaftState {
	stateAddr := (*uint32)(&r.state)
	return RaftState(atomic.LoadUint32((stateAddr)))
}

// 设置当前角色
func (r *raftState) setState(s RaftState) {
	stateAddr := (*uint32)(&r.state)
	atomic.StoreUint32(stateAddr, uint32(s))
}

// 获取当前任期
func (r *raftState) getCurrentTerm() uint64 {
	currentTerm := atomic.LoadUint64(&r.currentTerm)
	return currentTerm
}

// 设置当前任期
func (r *raftState) setCurrentTerm(term uint64) {
	atomic.StoreUint64(&r.currentTerm, term)
}

// 获取最后一次提交的索引号
func (r *raftState) getCommitIndex() uint32 {
	return atomic.LoadUint32(&r.commitIndex)
}

func (r *raftState) getCurrentCommitIndex() uint32 {
	index := atomic.LoadUint32(&r.commitIndex)
	// 因为uint32中负一表示为MAX_LOG_INDEX_NUM,所以不能简单+1
	if index == MAX_LOG_INDEX_NUM {
		index = 0
	} else {
		index++
	}
	return index
}

// 设置最后一次提交的索引号
func (r *raftState) setCommitIndex(index uint32) {
	atomic.StoreUint32(&r.commitIndex, index)
}

// 设置已投票给候选者中候选者最大任期号
func (r *raftState) setLatestVoteGrantedTerm(index uint64) {
	atomic.StoreUint64(&r.latestVoteGrantedTerm, index)
}

// 获取已投票给候选者中候选者最大任期号
func (r *raftState) getLatestVoteGrantedTerm() uint64 {
	return atomic.LoadUint64(&r.latestVoteGrantedTerm)
}

// 执行go routine
func (r *raftState) goFunc(f func()) {
	r.routinesGroup.Add(1)
	go func() {
		defer r.routinesGroup.Done()
		f()
	}()
}

// 等待关机
func (r *raftState) waitShutdown() {
	r.routinesGroup.Wait()
}

// 返回最后提交日志的索引号和任期
func (r *raftState) getLastEntry() (uint32, uint64) {
	r.lastLock.Lock()
	defer r.lastLock.Unlock()
	return r.lastLogIndex, r.lastLogTerm
}

// 返回当前日志的索引号和任期
func (r *raftState) getCurrentLogIndex() uint32 {
	r.lastLock.Lock()
	lastLogIndex := r.lastLogIndex
	r.lastLock.Unlock()

	// 因为uint32中负一表示为MAX_LOG_INDEX_NUM,所以不能简单+1
	if lastLogIndex == MAX_LOG_INDEX_NUM {
		lastLogIndex = 0
	} else {
		lastLogIndex++
	}
	return lastLogIndex
}

// 设置最后提交日志的索引号和任期
func (r *raftState) setLastEntry(lastLogIndex uint32, lastLogTerm uint64) {
	r.lastLock.Lock()
	defer r.lastLock.Unlock()
	r.lastLogIndex = lastLogIndex
	r.lastLogTerm = lastLogTerm
}
