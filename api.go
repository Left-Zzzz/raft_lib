package raftlib

import (
	"sync"
	"sync/atomic"
)

type Raft struct {
	raftState

	// leaderAddr: 当前集群中领导者网络地址
	leaderAddr ServerAddress
	// LeaderID: 当前集群中领导者id
	leaderID ServerID
	// 进行上面两项熟悉操作时用到的锁
	leaderLock sync.RWMutex

	// 集群节点数
	nodeNum uint64

	// 当前服务器ID
	localID ServerID
	// 当前服务器网络地址
	localAddr ServerAddress

	// rpc句柄
	rpc Rpc
}

// 获取领导者的网络地址
func (r *Raft) Leader() ServerAddress {
	r.leaderLock.RLock()
	defer r.leaderLock.RUnlock()
	leaderAddr := r.leaderAddr
	return leaderAddr
}

// 获取领导者的网络地址和id
func (r *Raft) LeaderWithID() (ServerAddress, ServerID) {
	r.leaderLock.RLock()
	defer r.leaderLock.RUnlock()
	leaderID := r.leaderID
	leaderAddr := r.leaderAddr
	return leaderAddr, leaderID
}

func (r *Raft) getNodeNum() uint64 {
	return atomic.LoadUint64(&r.nodeNum)
}

func (r *Raft) setNodeNum(num uint64) {
	atomic.StoreUint64(&r.nodeNum, num)
}
