package raftlib

import (
	"raft_lib/pb"
	"strconv"
	"time"
)

// run the main thread that handles leadership and RPC requests.
func (r *Raft) run() {
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
	heartbeatTimer := randomTimeout(HeartbeatTimeout)
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
			heartbeatTimer = randomTimeout(HeartbeatTimeout)

			// 如果在当轮计时中有收到rpc信息，则不算超时，循环继续
			lastContact := r.getLastContact()
			isContacted := lastContact.After(heartbeatStartTime)
			if isContacted {
				// 重置当前任期超时时间和计时开始时间
				heartbeatTimer = randomTimeout(HeartbeatTimeout)
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
	electionTimer := randomTimeout(ElectionTimeout)
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
					r.setState(Follower)
					continue
				}
				// 如果不是最新的，则continue
				continue
			}
			logDebug("recieve an vote agree from id:%s, term:%d.\n", vote.VoterID, vote.Term)
			voteGained++
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
	server_id, _ := strconv.Atoi(string(r.getLocalID()))
	if !isPrintRpcAppendEntryRequestCh[server_id] {
		logInfo("processRpcRequest(): rpcExecCommandRequestCh addr:%p, id:%d.\n", r.rpc.rpcCh.rpcExecCommandRequestCh, server_id)
		isPrintRpcAppendEntryRequestCh[server_id] = true
	}
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
		// do noting
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
	respCh := make(chan *pb.RequestVoteResponse, len(Servers))

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
				logError("failed to make requestVote RPC",
					"target", peer,
					"error", err,
					"term", req.Term)
				resp.Term = req.Term
				resp.VoteGranted = false
			}
			respCh <- resp
		})
	}

	// 遍历集群中其他节点，发送请求以获取投票
	logDebug("start vote. local.ID:%v", r.localID)
	for _, server := range Servers {
		// logDebug("server.ID:", server.ID, ". Suffrage:", server.Suffrage)
		if server.Suffrage == Voter {
			if server.ID == r.localID {
				logDebug("voting for self, term:%v, localID:%d", r.getCurrentTerm(), r.localID)
				// 自己给自己一票
				respCh <- &pb.RequestVoteResponse{
					Term:        req.Term,
					VoteGranted: true,
					VoterID:     string(r.localID),
				}
			} else {
				logDebug("asking for vote, term:%v, from:%v, address:%v", req.Term, server.ID, server.Address)
				go askPeer(server)
			}
		}
	}

	return respCh
}

func (r *Raft) sendHeartBeatLoop() {
	senHeartBeatInterval := HeartbeatTimeout / 2
	heartBeatTimer := time.After(senHeartBeatInterval)
	for r.getState() == Leader {
		<-heartBeatTimer
		logDebug("leader send heart beat.")
		prevLogIndex, prevLogTerm := r.raftState.getLastEntry()
		req := &pb.AppendEntryRequest{
			Ver:          &RPOTO_VER_APPEND_ENTRY_REQUEST,
			LeaderTerm:   r.getCurrentTerm(),
			LeaderID:     string(r.Leader()),
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			LogType:      uint32(HeartBeat),
			LeaderCommit: r.getCommitIndex(),
		}
		for _, server := range Servers {
			// 不用给自己发送RPC请求
			if server.ID == r.Leader() {
				continue
			}
			// TODO 处理RPC回复
			go r.rpc.AppendEntryRequest(server, req)
			logDebug("process AppendEntryResponse.")
		}
		heartBeatTimer = time.After(senHeartBeatInterval)
	}

}

// NoOp补丁，提交空日志
func (r *Raft) noOp() {
	// 直接复用execCommand
	req := &pb.ExecCommandRequest{
		Ver:     &PROTO_VER_EXEC_COMMAND_REQUEST,
		LogType: uint32(LogNoOp),
		Command: []byte(""),
	}
	resp := r.execCommand(req)
	if !resp.GetSuccess() {
		logWarn("r.execCommand(): failed. leaderAddress:%v, leaderPort:%v", resp.GetLeaderAddress(), resp.GetLeaderPort())
	}
}
