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
			// 收到拒绝投票，直接continue
			if !vote.GetVoteGranted() {
				logDebug("recieve an vote reject from id:%s, term:%d.\n", vote.VoterID, vote.Term)
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
	// TODO:发送当选leader通知
	r.sendHeartBeatLoop()
}

// 处理RPC请求
func (r *Raft) processRpcRequest() {
	server_id, _ := strconv.Atoi(string(r.getLocalID()))
	if !isPrintRpcAppendEntryRequestCh[server_id] {
		logDebug("processRpcRequest(): rpcAppendEntryRequestCh addr:%p.\n", r.rpc.rpcCh.rpcAppendEntryRequestCh)
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
		resp := r.appendEntryRpc(appendEntryRequest)
		r.rpc.rpcCh.rpcAppendEntryResponseCh <- resp
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
	lastIdx, lastTerm := r.getLastEntry()
	req := &pb.RequestVoteRequest{
		Term:        &curTerm,
		CandidateId: (*string)(&r.localID),
		LastLogIdx:  &lastIdx,
		LastLogTerm: &lastTerm,
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
					"term", *(req.Term))
				resp.Term = *(req.Term)
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
					Term:        *(req.Term),
					VoteGranted: true,
					VoterID:     string(r.localID),
				}
			} else {
				logDebug("asking for vote, term:%v, from:%v, address:%v", *req.Term, server.ID, server.Address)
				go askPeer(server)
			}
		}
	}

	return respCh
}

func (r *Raft) sendHeartBeatLoop() {
	logInfo("send heart beat.")
	req := &pb.AppendEntryRequest{
		Ver:        &RPOTO_VER_APPEND_ENTRY_REQUEST,
		LeaderTerm: r.getCurrentTerm(),
		LeaderID:   string(r.Leader()),
		LogType:    uint64(HeartBeat),
	}
	for _, server := range Servers {
		// TODO 处理RPC回复
		go r.rpc.AppendEntryRequest(server, req)
		logDebug("process AppendEntryResponse.")
	}
}
