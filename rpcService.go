package raftlib

import "raft_lib/pb"

// Raft层 requestVoteRPC 处理逻辑
func (r *Raft) requestVote(req *pb.RequestVoteRequest) (resp *pb.RequestVoteResponse) {
	resp = &pb.RequestVoteResponse{
		Ver:         &RPOTO_VER_REQUEST_VOTE_RESPONSE,
		Term:        r.getCurrentTerm(),
		VoteGranted: false,
	}

	// 如果本节点已有领导者，拒绝投票
	if leaderID := r.Leader(); leaderID != "" {
		return resp
	}

	// 如果follower任期比candidate大，拒绝投票
	if r.getCurrentTerm() > req.GetTerm() {
		return resp
	}

	// 如果follower日志比candidate日志新，拒绝投票
	lastLogIndex, lastLogTerm := r.getLastEntry()
	if lastLogIndex > *req.LastLogIdx {
		return resp
	}
	if lastLogIndex == *req.LastLogIdx && lastLogTerm > *req.LastLogTerm {
		return resp
	}

	// 如果已投票给候选者中候选者最大任期号不小于candidate任期号，拒绝投票
	if r.getLatestVoteGrantedTerm() >= *req.Term {
		return resp
	}

	// 满足条件，可以投票
	resp.VoteGranted = true
	// 记录已投票给候选者中候选者最大任期号
	r.setLatestVoteGrantedTerm(*req.Term)
	return resp
}

// 处理AppendEntryRpc请求
func (r *Raft) appendEntryRpc(req *pb.AppendEntryRequest) (resp *pb.AppendEntryResponse) {
	// 获取当前任期
	currentTerm := r.getCurrentTerm()
	// 如果是新leader，且新leader任期不小于本节点任期，则设置leaderId为新leader的Id
	rpcLeaderID := ServerID(req.GetLeaderID())
	rpcLeaderTerm := req.GetLeaderTerm()
	logDebug("rpcLeaderID: %s, rpcLeaderTerm:%d\n", string(rpcLeaderID), rpcLeaderTerm)
	if r.Leader() != rpcLeaderID && currentTerm <= rpcLeaderTerm {
		r.setLeaderID(rpcLeaderID)
		r.setCurrentTerm(rpcLeaderTerm)
	}

	resp = &pb.AppendEntryResponse{
		Term:    currentTerm,
		Success: false,
	}
	switch LogType(req.GetLogType()) {
	case HeartBeat:
		// TODO: 后续有问题再补充
		resp.Success = true
	case LogCommand:
		// TODO: 执行store log操作
	case LogNoOp:
		// TODO: No Op补丁，执行store log操作
	}
	r.rpc.rpcCh.rpcAppendEntryResponseCh <- resp
	return resp
}
