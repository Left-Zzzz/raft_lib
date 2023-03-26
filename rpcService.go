package raftlib

import (
	"encoding/json"
	"raft_lib/pb"
)

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

	// 将上一个日志的索引任期号添加到response中
	lastLogIndex, lastLogTerm := r.getLastEntry()
	resp.LastLogIdx = lastLogIndex
	resp.LastLogTerm = lastLogTerm

	// 如果follower日志比candidate日志新，拒绝投票
	if lastLogIndex > req.LastLogIdx {
		return resp
	}
	if lastLogIndex == req.LastLogIdx && lastLogTerm > req.LastLogTerm {
		return resp
	}

	// 如果已投票给候选者中候选者最大任期号不小于candidate任期号，拒绝投票
	if r.getLatestVoteGrantedTerm() >= req.Term {
		return resp
	}

	// 满足条件，可以投票
	resp.VoteGranted = true
	// 记录已投票给候选者中候选者最大任期号
	r.setLatestVoteGrantedTerm(req.Term)
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
	commitEntries := func() {
		r.goFunc(func() {
			leaderCommitIndex := req.GetLeaderCommit()
			if leaderCommitIndex == MAX_LOG_INDEX_NUM {
				// 没有日志提交，返回
				return
			}
			// 从map中获取应用日志的回调函数
			cbFunc, ok := r.storage.callBackFuncMap.Load(execCommandFuncName)
			if !ok {
				logError("callBackFunc load %v falied!", execCommandFuncName)
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
	if !r.storage.isIdxTermCorrect(req.GetPrevLogIndex(), req.GetPrevLogTerm()) {
		return resp
	}
	// 如果能对应，构造日志项
	logType := LogType(req.GetLogType())
	logEntry := req.GetEntry()
	// 获取应该被插入日志项的索引号
	currentLogIndex := genNextLogIndex(req.GetPrevLogIndex())
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
		err := r.storage.appendEntryEncoded(currentLogIndex, logEntry)
		if err != nil {
			logError("r.storage.appendEntry():%v", err)
			return resp
		}
		// 如果成功，将日志索引号和任期号缓存
		r.setLastEntry(genNextLogIndex(req.GetPrevLogIndex()), currentTerm)
		logDebug("r.setLastEntry(): idx: %v, term: %v", req.GetPrevLogIndex(), currentTerm)
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
	for _, server := range Servers {
		logInfo("server.ID:%v, leaderID:%v", server.ID, leaderID)
		if server.ID == leaderID {
			leaderAddress = string(server.Address)
			leaderPort = string(server.Port)
			logInfo("execCommand(), leaderAddress:%v, leaderPort:%v", string(server.Address), string(server.Port))
			break
		}
	}

	// 构造ExecCommandResponse
	resp := &pb.ExecCommandResponse{
		Ver:           &PROTO_VER_EXEC_COMMAND_RESPONSE,
		LeaderAddress: leaderAddress,
		LeaderPort:    leaderPort,
		Success:       false,
	}

	// 如果被请求的节点不是leader节点，返回leader所在地址
	if localID != leaderID {
		return resp
	}

	// 构造AppendEntry请求
	prevLogIndex, prvLogTerm := r.getLastEntry()
	logDebug("prevLogIndex:%v, prvLogTerm:%v", prevLogIndex, prvLogTerm)
	currentTerm := r.getCurrentTerm()
	appendEntryResponseCh := make(chan *pb.AppendEntryResponse)
	logEntryEncoded, _ := json.Marshal(&Log{
		Index:   prevLogIndex + 1,
		Term:    r.getCurrentTerm(),
		LogType: LogCommand,
		Data:    req.Command,
	})
	appendEntryReq := &pb.AppendEntryRequest{
		Ver:          &RPOTO_VER_APPEND_ENTRY_REQUEST,
		LeaderTerm:   r.getCurrentTerm(),
		LeaderID:     string(r.Leader()),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prvLogTerm,
		Entry:        logEntryEncoded,
		LeaderCommit: r.getCommitIndex(),
		LogType:      uint32(LogCommand),
	}
	// AppendEntryRPC通信流程
	askPeer := func(peer Server) {
		r.goFunc(func() {
			var isFirstTimeLoop bool = true
			var err error
			var resp *pb.AppendEntryResponse
			// 循环发送直到被接收回复为止
			for isFirstTimeLoop || err != nil {
				isFirstTimeLoop = false
				resp, err = r.rpc.AppendEntryRequest(peer, appendEntryReq)
			}
			appendEntryResponseCh <- resp
		})
	}
	// 创建协程发送请求
	for _, node := range Servers {
		// 忽略自己
		if node.ID == r.getLocalID() {
			continue
		}
		askPeer(node)
	}

	// 监听投票结果
	successCnt := 0
	leastSuccessRequired := r.getNodeNum()/2 + 1
	// 自己执行appendEntry操作
	localResp := r.appendEntry(appendEntryReq)
	// 如果leader节点都不能处理AppendEntry请求，返回错误
	if !localResp.GetSuccess() {
		return resp
	}
	// 更新上一个添加的日志项信息缓存
	r.setLastEntry(genNextLogIndex(prevLogIndex), currentTerm)
	logDebug("r.setLastEntry: idx:%v, term:%v", r.getCurrentLogIndex(), currentTerm)
	// 加上leader节点的一票
	successCnt++
	for {
		ch := <-appendEntryResponseCh
		// 如果收到成功响应消息
		if ch.GetSuccess() {
			successCnt++
			logDebug("execCommand(): recieve sucess response. success count:%v, success required:%v", successCnt, leastSuccessRequired)
			// 成功响应的节点超过一半
			if successCnt >= int(leastSuccessRequired) {
				// 从map中拿到回调函数
				cbFunc, ok := r.storage.callBackFuncMap.Load(execCommandFuncName)
				if !ok {
					logError("callBackFunc load %v falied!", execCommandFuncName)
					return resp
				}
				execCommandFunc, ok := cbFunc.(func([]byte) error)
				if !ok {
					logError("cbFunc transfer type (func([]byte) error) failed!")
					return resp
				}
				// TODO: 执行entry中的命令
				err := r.storage.commit(r.getCurrentCommitIndex(), execCommandFunc)
				if err != nil {
					logError("r.storage.commit(): %v, index:%v", err, r.getCurrentCommitIndex())
					return resp
				}
				// 打印输出已提交的日志项的信息
				debugLogEntry, err := r.storage.getEntry(r.getCurrentCommitIndex())
				if err != nil {
					logError("r.storage.getEntry(): %v", err)
					return resp
				}
				logDebug("r.storage.getEntry(%v): %v, command:%v", r.getCommitIndex(), debugLogEntry, string(debugLogEntry.Data))
				// 执行完成，更新CommitIndex
				r.setCommitIndex(r.getCurrentCommitIndex())

				// 回复成功响应
				logDebug("execCommand(): recieve over half success, repsonce true and exec command.")
				resp.Success = true
				return resp
			}
		}
	}
}
