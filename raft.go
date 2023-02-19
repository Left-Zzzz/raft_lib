package raftlib

import (
	"log"
	"raft_lib/pb"
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
	leaderAddr, leaderID := r.LeaderWithID()
	heartbeatTimer := randomTimeout(HeartbeatTimeout)
	log.Println("entering follower state", "follower", r, "leader-address", leaderAddr, "leader-id", leaderID)
	for r.getState() == Follower {
		select {
		// 处理RPC命令
		case <-heartbeatTimer:
			// 重置心跳计时器
			heartbeatTimer = randomTimeout(HeartbeatTimeout)

			// TODO: 如果在当轮计时中有收到rpc信息，则不算超时，循环继续
			isContacted := true
			if isContacted {
				continue
			}

			// 超时处理
			// 删除leader相关信息
			r.setLeader("", "")
			// 将自身角色设置为candidate
			r.setState(Candidate)
		}
	}

}

func (r *Raft) runCandidate() {
	// 缓存任期
	curTerm := r.getCurrentTerm() + 1

	voteCh := r.election()

	// 设置选举计时器
	electionTimer := randomTimeout(ElectionTimeout)
	voteGained := 0
	leastVotesRequired := r.getNodeNum()/2 + 1

	for r.getState() == Candidate {

	}
}

func (r *Raft) runLeader() {

}

func (r *Raft) setLeader(leaderAddr ServerAddress, leaderID ServerID) {
	r.leaderLock.Lock()
	defer r.leaderLock.Unlock()
	r.leaderAddr = leaderAddr
	r.leaderID = leaderID
}

type voteResult struct {
	pb.VoteResponse
	voterID ServerID
}

func (r *Raft) electSelf() <-chan *voteResult {
	// 新建response通道
	respCh := make(chan *voteResult, len(Servers))

	// 节点任期号+1
	r.setCurrentTerm(r.getCurrentTerm() + 1)

	// 构造RPC请求message
	curTerm := r.getCurrentTerm()
	lastIdx, lastTerm := r.getLastEntry()
	req := &pb.VoteRequest{
		Term:        &curTerm,
		CandidateId: (*string)(&r.localID),
		LastLogIdx:  &lastIdx,
		LastLogTerm: &lastTerm,
		//LeadershipTransfer: r.candidateFromLeadershipTransfer,
	}

	// TODO: 构造选举请求函数，构造responseMessage
	askPeer := func(peer Server) {
		r.goFunc(func() {
			// defer metrics.MeasureSince([]string{"raft", "candidate", "electSelf"}, time.Now())
			resp := &voteResult{voterID: peer.ID}
			// 发送RPC请求
			// err := r.trans.RequestVote(peer.ID, peer.Address, req, &resp.RequestVoteResponse)
			err := r.rpc.RequestVote(peer.ID, peer.Address, req, &resp.VoteResponse)
			if err != nil {
				log.Println("[Error] failed to make requestVote RPC",
					"target", peer,
					"error", err,
					"term", req.Term)
				resp.Term = *req.Term
				resp.VoteGranted = false
			}
			respCh <- resp
		})
	}

	// 遍历集群中其他节点，发送请求以获取投票
	for _, server := range Servers {
		if server.Suffrage == Voter {
			if server.ID == r.localID {
				log.Println("[Debug] voting for self", "term", r.getCurrentTerm(), "id", r.localID)
				// 自己给自己一票
				/*
					respCh <- &voteResult{
						RequestVoteResponse: RequestVoteResponse{
							RPCHeader: r.getRPCHeader(),
							Term:      req.Term,
							Granted:   true,
						},
						voterID: r.localID,
					}
				*/
			} else {
				log.Println("[Debug] asking for vote", "term", req.Term, "from", server.ID, "address", server.Address)
				askPeer(server)
			}
		}
	}

	return respCh
}
