package raftlib

import (
	"encoding/json"
	"raftlib/pb"
)

// Raft层 requestVoteRPC 处理逻辑
func (r *Raft) requestVote(req *pb.RequestVoteRequest) (resp *pb.RequestVoteResponse) {
	ver := RPOTO_VER_REQUEST_VOTE_RESPONSE
	lastLogIndex, lastLogTerm := r.getLastEntry()
	resp = &pb.RequestVoteResponse{
		Ver:         &ver,
		VoterID:     string(r.localID),
		Term:        r.getCurrentTerm(),
		LastLogIdx:  lastLogIndex,
		LastLogTerm: lastLogTerm,
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

	// 安全性补丁：判断日志是否是最新的，如果不是最新的话，candidate节点回退至follower
	// 这样做法是保证日志是安全的，不会被日志少的节点当选leader，进而造成已提交日志被覆盖情况

	// 如果follower日志比candidate日志新，拒绝投票
	logDebug("lastLogIndex:%v, req.LastLogIdx:%v", lastLogIndex, req.GetLastLogIdx())
	logDebug("lastLogTerm: %v, req.LastLogTerm: %v", lastLogTerm, req.GetLastLogTerm())
	if genActualLogIndex(lastLogIndex) > genActualLogIndex(req.GetLastLogIdx()) {
		return resp
	}
	if lastLogIndex == req.GetLastLogIdx() && lastLogTerm > req.GetLastLogTerm() {
		return resp
	}

	// 如果已投票给候选者中候选者最大任期号不小于candidate任期号，拒绝投票
	if r.getLatestVoteGrantedTerm() >= req.GetTerm() {
		return resp
	}

	// 满足条件，可以投票
	resp.VoteGranted = true
	// 记录已投票给候选者中候选者最大任期号
	r.setLatestVoteGrantedTerm(req.GetTerm())
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

	// 当前节点同步提交leader已提交的logEntry
	commitEntries := func() {
		r.goFunc(func() {
			leaderCommitIndex := req.GetLeaderCommit()
			logDebug("commitEntries(): leaderCommitIndex:%v", leaderCommitIndex)
			if leaderCommitIndex == MAX_LOG_INDEX_NUM {
				// 没有日志提交，返回
				return
			}
			// 从map中获取应用日志的回调函数
			cbFunc, ok := r.storage.callBackFuncMap.Load(EXEC_COMMAND_FUNC_NAME)
			if !ok {
				logError("callBackFunc load %v falied!", EXEC_COMMAND_FUNC_NAME)
				return
			}
			execCommandFunc, ok := cbFunc.(func([]byte) error)
			if !ok {
				logError("cbFunc transfer type (func([]byte) error) failed!")
				return
			}
			// 应用日志
			commitIndex := req.GetLeaderCommit()
			logDebug("r.storage.batchCommit(), commitIndex:%v", commitIndex)
			r.storage.batchCommit(commitIndex, execCommandFunc)
			// 更新缓存
			r.setCommitIndex(r.storage.getCommitIndex())
		})
	}
	commitEntries()

	// 判断AppendEntryRpc请求日志是否对应，如果日志不能对应，返回失败
	reqPrevLogIndex := req.GetPrevLogIndex()
	reqPrevLogTerm := req.GetPrevLogTerm()
	logDebug("req.GetPrevLogIndex():%v, req.GetPrevLogTerm():%v", reqPrevLogIndex, reqPrevLogTerm)
	// Debug对比
	nodeEntry, err := r.storage.getEntry(reqPrevLogIndex)
	if err != nil {
		logDebug("r.storage.getEntry():%v", err)
	} else {
		logDebug("nodeEntry.Index:%v, nodeEntry.LogTerm:%v", nodeEntry.Index, nodeEntry.Term)
	}

	// 判断上一个日志项是否匹配
	if !r.storage.isIdxTermCorrect(reqPrevLogIndex, reqPrevLogTerm) {
		logDebug("appendEntry(): log entry not match! req.GetPrevLogIndex(): %v,req.GetPrevLogTerm(): %v",
			req.GetPrevLogIndex(), req.GetPrevLogTerm())
		return resp
	}
	// 如果能对应，构造日志项
	logType := LogType(req.GetLogType())
	logEntryEncoded := req.GetEntry()
	// 获取应该被插入日志项的索引号
	logEntry := &Log{}
	err = json.Unmarshal(logEntryEncoded, logEntry)
	if err != nil {
		logWarn("json.Unmarshal(req.GetEntry(), logEntry):%v", err)
	}
	switch logType {
	case HeartBeat:
		// 如果是心跳包，不用执行日志项追加复制操作
		resp.Success = true
	case LogCommand:
		// 执行store log操作，因为是一样的操作，合并操作
		fallthrough
	case LogNoOp:
		// No Op补丁，执行store log操作
		logDebug("call r.storage.appendEntry()")
		err := r.storage.appendEntryEncoded(logEntry.Index, logEntryEncoded)
		if err != nil {
			logError("r.storage.appendEntry():%v", err)
			return resp
		}
		// 如果成功，将日志索引号和任期号缓存
		r.setLastEntry(logEntry.Index, logEntry.Term)
		logDebug("r.setLastEntry(): idx: %v, term: %v", logEntry.Index, logEntry.Term)
		resp.Success = true
	default:
		logWarn("Unkonwn logType.")
	}
	return resp
}

func (r *Raft) execCommand(req *pb.ExecCommandRequest) *pb.ExecCommandResponse {
	// 获取leader的address和port
	localID := r.getLocalID()
	leaderID := r.Leader()
	leaderAddress := ""
	leaderPort := ""
	for _, server := range r.config.Servers {
		logInfo("server.ID:%v, leaderID:%v", server.ID, leaderID)
		if server.ID == leaderID {
			leaderAddress = string(server.Address)
			leaderPort = string(server.Port)
			logInfo("execCommand(), leaderAddress:%v, leaderPort:%v", string(server.Address), string(server.Port))
			break
		}
	}

	// 构造ExecCommandResponse
	ver := PROTO_VER_EXEC_COMMAND_RESPONSE
	resp := &pb.ExecCommandResponse{
		Ver:           &ver,
		LeaderID:      string(leaderID),
		LeaderAddress: leaderAddress,
		LeaderPort:    leaderPort,
		Success:       false,
	}

	// 如果被请求的节点不是leader节点，返回leader所在地址
	if localID != leaderID {
		logDebug("localID(%v) != leaderID(%v), localAddr:%v, localPort:%v, return false.", localID, leaderID, r.localAddr, r.localPort)
		return resp
	}

	appendEntryResponseCh := make(chan *pb.AppendEntryResponse)
	currentTerm := r.getCurrentTerm()

	// 构造AppendEntryRPC请求，因为是向其他节点发送请求，所以请求的日志leader节点都有
	genAppendEntryRequest := func(logIndex uint32) *pb.AppendEntryRequest {
		logEntry, err := r.storage.getEntry(logIndex)
		if err != nil {
			logError("genAppendEntryRequest:r.storage.getEntry(): %v", err)
		}
		logType := logEntry.LogType
		logEntryEncoded, err := json.Marshal(logEntry)
		if err != nil {
			logError("genAppendEntryRequest:json.Marshal(): %v", err)
		}
		prevLogIndex := genPrevLogIndex(logIndex)
		prevLogTerm := uint64(0)
		if prevLogIndex != MAX_LOG_INDEX_NUM {
			prevLogEntry, err := r.storage.getEntry(prevLogIndex)
			if err != nil {
				logError("genAppendEntryRequest:r.storage.getEntry(): %v", err)
			} else {
				logDebug("prevLogEntry:%v", prevLogEntry)
			}
			prevLogTerm = prevLogEntry.Term
		}
		logDebug("genAppendEntryRequest(): prevLogIndex: %v, prevLogTerm: %v", prevLogIndex, prevLogTerm)
		ver := RPOTO_VER_APPEND_ENTRY_REQUEST
		return &pb.AppendEntryRequest{
			Ver:          &ver,
			LeaderTerm:   currentTerm,
			LeaderID:     string(r.Leader()),
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entry:        logEntryEncoded,
			LeaderCommit: r.getCommitIndex(),
			LogType:      uint32(logType),
		}
	}

	// AppendEntryRPC通信流程
	var askPeer func(Server, uint32) = func(peer Server, logIndex uint32) {
		r.goFunc(func() {
			var err error
			var resp *pb.AppendEntryResponse
			// 如果遇到日志缺漏时，回退以下三个参数
			currentLogIndex := logIndex
			for {
				// 构造AppendEntryRequest
				tempAppendEntryReq := genAppendEntryRequest(currentLogIndex)
				logDebug("askPeer: tempAppendEntryReq: %v", tempAppendEntryReq)
				// 循环发送直到被接收回复为止
				for {
					resp, err = r.rpc.AppendEntryRequest(peer, tempAppendEntryReq)
					if err != nil {
						logDebug("r.rpc.AppendEntryRequest(): %v", err)
					} else {
						break
					}
				}
				// 如果resp.Success为true, 继续提交下一个日志
				if resp.GetSuccess() {
					logDebug("askPeer: response success, currentLogIndex:%v", currentLogIndex)
					// 如果currentPrevLogIndex == prevLogIndex，结束
					if currentLogIndex == logIndex {
						logDebug("askPeer: currentLogIndex == logIndex, break.")
						break
					}
					currentLogIndex = genNextLogIndex(currentLogIndex)
				} else {
					logDebug("askPeer: response false, currentLogIndex:%v", currentLogIndex)
					// 如果resp.Success为false, 回溯提交上一个日志
					// 终止回溯条件
					if currentLogIndex == MAX_LOG_INDEX_NUM {
						logError("askPeer: prevLogIndex == MAX_LOG_INDEX_NUM")
						break
					}
					currentLogIndex = genPrevLogIndex(currentLogIndex)
					logDebug("log entry not match, continue askPeer(currentLogIndex:%v)", currentLogIndex)
				}
			}
			appendEntryResponseCh <- resp

		})
	}

	// 构造AppendEntry请求
	prevLogIndex, prevLogTerm := r.getLastEntry()
	currentLogIndex := genNextLogIndex(prevLogIndex)
	logDebug("prevLogIndex:%v, prvLogTerm:%v", prevLogIndex, prevLogTerm)

	logEntryEncoded, err := json.Marshal(&Log{
		Index:   currentLogIndex,
		Term:    currentTerm,
		LogType: LogType(req.GetLogType()),
		Data:    req.GetCommand(),
	})
	if err != nil {
		logError(" .json.Marshal():%v", err)
	}
	// 监听投票结果
	successCnt := 0
	leastSuccessRequired := r.getNodeNum()/2 + 1
	// leader节点自己执行appendEntry操作
	ver = RPOTO_VER_APPEND_ENTRY_REQUEST
	appendEntryReq := &pb.AppendEntryRequest{
		Ver:          &ver,
		LeaderTerm:   currentTerm,
		LeaderID:     string(r.Leader()),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entry:        logEntryEncoded,
		LeaderCommit: r.getCommitIndex(),
		LogType:      req.GetLogType(),
	}
	logDebug("leader self appendEntryReq:%v", appendEntryReq)
	localResp := r.appendEntry(appendEntryReq)
	// 如果leader节点都不能处理AppendEntry请求，返回错误
	if !localResp.GetSuccess() {
		logWarn("leader node can not process AppendEntry Request.")
		return resp
	}
	// 加上leader节点的一票
	successCnt++

	// 创建协程向其他节点发送请求
	for _, node := range r.config.Servers {
		// 忽略自己
		if node.ID == r.getLocalID() {
			continue
		}
		askPeer(node, currentLogIndex)
	}

	for {
		ch := <-appendEntryResponseCh
		// 如果收到成功响应消息
		if !ch.GetSuccess() {
			continue
		}
		successCnt++
		logDebug("execCommand(): recieve sucess response. success count:%v, success required:%v", successCnt, leastSuccessRequired)
		if successCnt < int(leastSuccessRequired) {
			continue
		}
		// 成功响应的节点超过一半
		// 从map中拿到回调函数
		cbFunc, ok := r.storage.callBackFuncMap.Load(EXEC_COMMAND_FUNC_NAME)
		if !ok {
			logError("callBackFunc load %v falied!", EXEC_COMMAND_FUNC_NAME)
			return resp
		}
		execCommandFunc, ok := cbFunc.(func([]byte) error)
		if !ok {
			logError("cbFunc transfer type (func([]byte) error) failed!")
			return resp
		}
		// 执行entry中的命令, 将未提交的entry一并提交
		err := r.storage.batchCommit(currentLogIndex, execCommandFunc)
		if err != nil {
			logError("r.storage.batchCommit(): %v, index:%v", err, currentLogIndex)
			return resp
		}
		// 执行完成，更新CommitIndex
		r.setCommitIndex(currentLogIndex)

		// 回复成功响应
		logDebug("execCommand(): recieve over half success, repsonce true and exec command.")
		resp.Success = true
		return resp
	}
}
