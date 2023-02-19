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

	// 当前任期内收到选票的candidateId，如果没有投给任何候选人则为空
	votedFor uint64

	// 最大日志提交索引号，初始值0，单增
	commitIndex uint64

	// 已经被应用到状态机的最高日志条目索引号，初始值0，单增
	lastApplied uint64

	// 保护一下两个参数的读写操作
	lastLock sync.Mutex

	// 缓存最新log Index/Term
	lastLogIndex uint64
	lastLogTerm  uint64

	// 记录正在运行的goroutines数量
	routinesGroup sync.WaitGroup

	// 当前状态机角色
	state RaftState
}

func (r *raftState) getState() RaftState {
	stateAddr := (*uint32)(&r.state)
	return RaftState(atomic.LoadUint32((stateAddr)))
}

func (r *raftState) setState(s RaftState) {
	stateAddr := (*uint32)(&r.state)
	atomic.StoreUint32(stateAddr, uint32(s))
}

func (r *raftState) getCurrentTerm() uint64 {
	return atomic.LoadUint64(&r.currentTerm)
}

func (r *raftState) setCurrentTerm(term uint64) {
	atomic.StoreUint64(&r.currentTerm, term)
}

func (r *raftState) getCommitIndex() uint64 {
	return atomic.LoadUint64(&r.commitIndex)
}

func (r *raftState) setCommitIndex(index uint64) {
	atomic.StoreUint64(&r.commitIndex, index)
}

func (r *raftState) getLastApplied() uint64 {
	return atomic.LoadUint64(&r.lastApplied)
}

func (r *raftState) setLastApplied(index uint64) {
	atomic.StoreUint64(&r.lastApplied, index)
}

func (r *raftState) goFunc(f func()) {
	r.routinesGroup.Add(1)
	go func() {
		defer r.routinesGroup.Done()
		f()
	}()
}

func (r *raftState) waitShutdown() {
	r.routinesGroup.Wait()
}

// getLastEntry returns the last index and term in stable storage.
// Either from the last log or from the last snapshot.
func (r *raftState) getLastEntry() (uint64, uint64) {
	r.lastLock.Lock()
	defer r.lastLock.Unlock()
	return r.lastLogIndex, r.lastLogTerm
}
