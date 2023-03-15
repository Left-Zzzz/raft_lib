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
func (r *Raft) appendEntry(req *pb.AppendEntryRequest) (resp *pb.AppendEntryResponse) {
	// 获取当前任期
	currentTerm := r.getCurrentTerm()
	// 构造AppendEntryResponse
	resp = &pb.AppendEntryResponse{
		Term:    currentTerm,
		Success: false,
	}
	// 如果是新leader
	rpcLeaderID := ServerID(req.GetLeaderID())
	rpcLeaderTerm := req.GetLeaderTerm()
	if r.Leader() != rpcLeaderID {
		// 如果新leader任期不小于本节点任期，则设置leaderId为新leader的Id
		if currentTerm <= rpcLeaderTerm {
			r.setLeaderID(rpcLeaderID)
			r.setCurrentTerm(rpcLeaderTerm)
			logDebug("rpcLeaderID: %s, rpcLeaderTerm:%d\n", string(rpcLeaderID), rpcLeaderTerm)
		} else {
			// 新leader任期小于本节点任期，拒绝请求
			return resp
		}
	}

	// TODO: 当前节点同步提交leader已提交的logEntry
	prevLogIndex, prevLogTerm := r.getLastEntry()
	commitEntries := func() {
		r.goFunc(func() {
			leaderCommitIndex := req.GetLeaderCommit()
			if leaderCommitIndex == MAX_LOG_INDEX_NUM {
				// 没有日志提交，返回
				return
			}
			lastCommitIndex := r.getCommitIndex()
			var curCommitIndex uint32
			// 计算当前待提交日志的下标
			if lastCommitIndex == MAX_LOG_INDEX_NUM {
				curCommitIndex = 0
			} else {
				curCommitIndex = lastCommitIndex + 1
			}

			for index := curCommitIndex; index <= leaderCommitIndex; index++ {
				err := r.storage.commit(index, execCommandfunc)
				if err != nil {
					logError("r.storage.commit():%v", err)
					break
				}
				// 更新最新应用的日志号
				r.setCommitIndex(index)
			}
		})
	}
	commitEntries()

	// 判断AppendEntryRpc请求日志是否对应，如果日志不能对应，返回失败
	if prevLogIndex != req.GetPrevLogIndex() || prevLogTerm != req.GetPrevLogTerm() {
		return resp
	}
	// 如果能对应，构造日志项
	logType := LogType(req.GetLogType())
	data := req.GetEntry()
	currentLogIndex := r.getCurrentLogIndex()
	logEntry := createLogEntry(currentLogIndex, logType, currentTerm, data)
	switch logType {
	case HeartBeat:
		// TODO: 后续有问题再补充
		resp.Success = true
	case LogCommand:
		// TODO: 执行store log操作，因为是一样的操作，合并操作
		fallthrough
	case LogNoOp:
		// TODO: No Op补丁，执行store log操作
		logDebug("call r.storage.appendEntry()")
		err := r.storage.appendEntry(logEntry)
		if err != nil {
			logError("r.storage.appendEntry():%v", err)
			return resp
		}
		// 如果成功，将日志索引号和任期号缓存
		r.setLastEntry(prevLogIndex+1, prevLogTerm+1)
		resp.Success = true
	default:
		logWarn("Unkonwn logType.")
	}
	return resp
}
