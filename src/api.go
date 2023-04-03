package raftlib

import (
	"raftlib/pb"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

type Raft struct {
	*raftState

	// LeaderID: 当前集群中领导者id
	leaderID ServerID
	// 进行上面操作时用到的锁
	leaderLock sync.RWMutex

	// 集群节点数
	nodeNum uint64

	// 当前服务器ID
	localID ServerID
	// 当前服务器网络地址
	localAddr ServerAddress
	// 当前服务器端口
	localPort ServerPort

	// log存储介质
	storage LogStorage

	// 最近一次RPC时间
	lastContact     time.Time
	lastContactLock sync.RWMutex

	// rpc句柄
	rpc *Rpc

	// 配置文件
	config *Config
}

// 获取领导者的Adress
func (r *Raft) Leader() ServerID {
	r.leaderLock.RLock()
	defer r.leaderLock.RUnlock()
	leaderID := r.leaderID
	return leaderID
}

func (r *Raft) getNodeNum() uint64 {
	return atomic.LoadUint64(&r.nodeNum)
}

func (r *Raft) setNodeNum(num uint64) {
	atomic.StoreUint64(&r.nodeNum, num)
}

func (r *Raft) setLastContact() {
	r.lastContactLock.Lock()
	defer r.lastContactLock.Unlock()
	r.lastContact = time.Now()
}

func (r *Raft) getLastContact() (lastContact time.Time) {
	r.lastContactLock.RLock()
	defer r.lastContactLock.RUnlock()
	return r.lastContact
}

func (r *Raft) getLocalID() ServerID {
	return r.localID
}

func (r *Raft) setLocalID(localID ServerID) {
	r.localID = localID
}

func CreateRaft(config *Config) *Raft {
	server := config.Localserver
	// 构造Raft
	r := &Raft{
		localID:   server.ID,
		localAddr: server.Address,
		localPort: server.Port,
		nodeNum:   uint64(len(config.Servers)),
		storage:   *createLogStorage(),
		raftState: &raftState{
			commitIndex:  MAX_LOG_INDEX_NUM,
			lastLogIndex: MAX_LOG_INDEX_NUM,
		},
		config: config,
	}
	debug = config.DEBUG
	r.setLeader(Server{})
	// 构造Rpc结构体
	rpc := &Rpc{
		clientConns: make(map[ServerID]*grpc.ClientConn),
		clients:     make(map[ServerID]pb.RpcServiceClient),
	}
	r.rpc = rpc

	// 构造gRPC Server
	var err any
	r.rpc.server, err = r.rpc.createRpcServer(r.config)
	go r.rpc.runRpcServer(r.rpc.server, server.Port)
	if err != nil {
		panic(err)
	}

	logDebug("CreateRaft(): rpcCh:%v", r.rpc.rpcCh)
	return r
}

func (r *Raft) RegisterCallBackFunc(f func([]byte) error) {
	logDebug("Raft.RegisterCallBackFunc(): localAddr:%v, localPort:%v", r.localAddr, r.localPort)
	r.storage.callBackFuncMap.Store(EXEC_COMMAND_FUNC_NAME, f)
}
