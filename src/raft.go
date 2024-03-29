package raftlib

import (
	"raftlib/pb"
	"time"
)

// Raft运行主流程
func (r *Raft) Run() {
	for {
		switch r.getState() {
		case Follower:
			r.runFollower()
		case Candidate:
			r.runCandidate()
		case Leader:
			r.runLeader()
		}
	}
}

func (r *Raft) runFollower() {
	leaderID := r.Leader()
	heartbeatTimer := randomTimeout(r.config.HeartbeatTimeout)
	heartbeatStartTime := time.Now()
	// 输出raft节点元信息
	logInfo("entering follower state, id: %s, term: %d, leader_id: %s\n",
		r.getLocalID(),
		r.getCurrentTerm(),
		leaderID)

	for r.getState() == Follower {
		// 处理RPC命令
		r.processRpcRequest()
		select {
		case <-heartbeatTimer:
			// 重置心跳计时器
			heartbeatTimer = randomTimeout(r.config.HeartbeatTimeout)

			// 如果在当轮计时中有收到rpc信息，则不算超时，循环继续
			lastContact := r.getLastContact()
			isContacted := lastContact.After(heartbeatStartTime)
			if isContacted {
				// 重置当前任期超时时间和计时开始时间
				heartbeatTimer = randomTimeout(r.config.HeartbeatTimeout)
				heartbeatStartTime = time.Now()
				continue
			}

			// 超时处理
			// 删除leader相关信息
			r.setLeader(Server{})
			// 将自身角色设置为candidate
			r.setState(Candidate)
		default:
			// do noting.
		}
	}

}

func (r *Raft) runCandidate() {
	// 缓存任期
	curTerm := r.getCurrentTerm() + 1
	// 输出raft节点元信息
	logInfo("entering candidate state, node: %s, term: %d\n", string(r.getLocalID()), int(curTerm))

	// 发起选举
	voteRespCh := r.electSelf()

	// 设置选举计时器
	electionTimer := randomTimeout(r.config.ElectionTimeout)
	voteGained := 0
	leastVotesRequired := r.getNodeNum()/2 + 1

	for r.getState() == Candidate {
		// 处理rpc命令
		r.processRpcRequest()
		select {
		case <-electionTimer:
			r.setState(Follower)
		case vote := <-voteRespCh:
			// 如果遇到任期更大的节点，回退当前节点角色至follower
			if vote.GetTerm() > curTerm {
				r.setState(Follower)
			}

			// 安全性补丁：判断日志是否是最新的，如果不是最新的话，candidate节点回退至follower
			// 这样做法是保证日志是安全的，不会被日志少的节点当选leader，进而造成已提交日志被覆盖情况

			// 收到拒绝投票，查看对方日志是否是最新的
			if !vote.GetVoteGranted() {
				logDebug("recieve an vote reject from id:%s, term:%d.\n", vote.VoterID, vote.Term)
				// 如果是最新的，则raftState回退至follower
				lastLogIdx, lastLogTerm := r.getLastEntry()
				if lastLogIdx < vote.GetLastLogIdx() || lastLogTerm < vote.GetLastLogTerm() {
					logDebug("recieve an newer log node from id:%s, term:%d.\n", vote.VoterID, vote.Term)
					r.setState(Follower)
					break
				}
				// 如果不是最新的，则continue
				continue
			}
			logDebug("recieve an vote agree from id:%s, term:%d.\n", vote.VoterID, vote.Term)
			voteGained++
			logDebug("voteGained:%v, leastVotesRequired:%v.\n", voteGained, leastVotesRequired)
			if leastVotesRequired <= uint64(voteGained) {
				// 设置角色为leader
				r.setState(Leader)
				// 设置leader为本机
				r.setLeader(Server{
					Address: r.localAddr,
					Port:    r.localPort,
					ID:      r.localID,
				})
			}

		default:
			// do nothing.
		}

	}
}

func (r *Raft) runLeader() {
	logInfo("entering leader state, node: %s, term: %d\n", string(r.getLocalID()), int(r.getCurrentTerm()))
	// 安全性补丁：NoOp补丁，leader提交非自身任期的日志是十分危险的，会导致已提交日志被覆盖
	// 所以leader节点只能提交自身任期的日志，而NoOp补丁既可以提交自身日志，又能将旧日志安全提交
	r.noOp()

	// 发送当选leader通知
	go r.sendHeartBeatLoop()
	for r.getState() == Leader {
		// 处理RPC请求
		r.processRpcRequest()
	}
}

// 处理RPC请求
func (r *Raft) processRpcRequest() {
	isContacted := true
	select {
	// 处理RequestVoteRequest RPC
	case requestVoteRequest := <-r.rpc.rpcCh.rpcRequestVoteRequestCh:
		logDebug("recieve an message from rpcRequestVoteRequestCh.\n")
		resp := r.requestVote(requestVoteRequest)
		r.rpc.rpcCh.rpcRequestVoteResponseCh <- resp
	// 处理AppendEntry RPC
	case appendEntryRequest := <-r.rpc.rpcCh.rpcAppendEntryRequestCh:
		logDebug("recieve an message from rpcAppendEntryRequestCh.\n")
		resp := r.appendEntry(appendEntryRequest)
		r.rpc.rpcCh.rpcAppendEntryResponseCh <- resp
	// 处理execCommand RPC
	case execCommandRequest := <-r.rpc.rpcCh.rpcExecCommandRequestCh:
		logInfo("recieve an message from rpcExecCommandRequestCh.\n")
		resp := r.execCommand(execCommandRequest)
		logInfo("process an message from rpcExecCommandRequestCh:%v.\n", resp)
		r.rpc.rpcCh.rpcExecCommandResponseCh <- resp
	default:
		isContacted = false
	}

	// 如果进行了RPC通信，将通信时间设置为最新通信时间
	if isContacted {
		r.lastContact = time.Now()
	}
}

func (r *Raft) setLeader(server Server) {
	r.leaderLock.Lock()
	defer r.leaderLock.Unlock()
	r.leaderID = server.ID
}

func (r *Raft) setLeaderID(id ServerID) {
	r.leaderLock.Lock()
	defer r.leaderLock.Unlock()
	r.leaderID = id
}

func (r *Raft) electSelf() <-chan *pb.RequestVoteResponse {
	// 新建response通道
	respCh := make(chan *pb.RequestVoteResponse, len(r.config.Servers))

	// 节点任期号+1
	r.setCurrentTerm(r.getCurrentTerm() + 1)

	// 构造RPC请求message
	curTerm := r.getCurrentTerm()
	lastLogIdx, lastLogTerm := r.getLastEntry()
	req := &pb.RequestVoteRequest{
		Term:        curTerm,
		CandidateId: string(r.localID),
		LastLogIdx:  lastLogIdx,
		LastLogTerm: lastLogTerm,
	}

	// 构造选举请求函数，构造responseMessage
	askPeer := func(peer Server) {
		r.goFunc(func() {
			// 发送RPC请求
			resp, err := r.rpc.RequestVoteRequest(peer, req)
			if err != nil {
				logDebug("resp:%p", resp)
				logError("failed to make requestVote RPC",
					"target", peer,
					"error", err,
					"term", req.GetTerm())
				resp.Term = req.GetTerm()
				resp.VoteGranted = false
			}
			respCh <- resp
		})
	}

	// 遍历集群中其他节点，发送请求以获取投票
	logDebug("start vote. local.ID:%v", r.localID)
	for _, server := range r.config.Servers {
		// logDebug("server.ID:", server.ID, ". Suffrage:", server.Suffrage)
		if server.Suffrage == Voter {
			if server.ID == r.localID {
				logDebug("voting for self, term:%v, localID:%d", r.getCurrentTerm(), r.localID)
				// 自己给自己一票
				respCh <- &pb.RequestVoteResponse{
					Term:        req.GetTerm(),
					VoteGranted: true,
					VoterID:     string(r.localID),
				}
			} else {
				logDebug("asking for vote, term:%v, from:%v, address:%v", req.GetTerm(), server.ID, server.Address)
				go askPeer(server)
			}
		}
	}

	return respCh
}

func (r *Raft) sendHeartBeatLoop() {
	senHeartBeatInterval := r.config.HeartbeatTimeout / 2
	heartBeatTimer := time.After(senHeartBeatInterval)
	respCh := make(chan *pb.AppendEntryResponse, 10)
	// 构造选举请求函数，构造responseMessage
	askPeer := func(peer Server, req *pb.AppendEntryRequest) {
		r.goFunc(func() {
			// 发送RPC请求
			resp, err := r.rpc.AppendEntryRequest(peer, req)
			if err != nil {
				logDebug("resp:%p", resp)
				return
			}
			respCh <- resp
		})
	}
	for r.getState() == Leader {
		<-heartBeatTimer
		logDebug("leader send heart beat.")
		prevLogIndex, prevLogTerm := r.raftState.getLastEntry()
		ver := RPOTO_VER_APPEND_ENTRY_REQUEST
		req := &pb.AppendEntryRequest{
			Ver:          &ver,
			LeaderTerm:   r.getCurrentTerm(),
			LeaderID:     string(r.Leader()),
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			LogType:      uint32(HeartBeat),
			LeaderCommit: r.getCommitIndex(),
		}
		for _, server := range r.config.Servers {
			// 不用给自己发送RPC请求
			if server.ID == r.Leader() {
				continue
			}
			// 处理RPC回复
			askPeer(server, req)
			logDebug("process AppendEntryResponse.")
		}

		// 当appendEntry请求收到false回复时，执行NoOp补丁，让日志落后的节点更新日志
		select {
		case resp := <-respCh:
			if !resp.GetSuccess() {
				logInfo("noOp")
				r.noOp()
			}
		default:
			// do nothing.
		}

		heartBeatTimer = time.After(senHeartBeatInterval)
	}
}

// NoOp补丁，提交空日志
func (r *Raft) noOp() {
	// 最多只能有一个NoOp执行体
	ok := r.noOpLock.TryLock()
	if !ok {
		return
	}
	defer r.noOpLock.Unlock()
	// 直接复用execCommand
	ver := PROTO_VER_EXEC_COMMAND_REQUEST
	req := &pb.ExecCommandRequest{
		Ver:     &ver,
		LogType: uint32(LogNoOp),
		Command: []byte(""),
	}
	resp := r.execCommand(req)
	if !resp.GetSuccess() {
		logWarn("r.execCommand(): failed. leaderAddress:%v, leaderPort:%v", resp.GetLeaderAddress(), resp.GetLeaderPort())
	}
}
