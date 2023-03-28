package raftlib

import (
	"encoding/json"
	"log"
	"raft_lib/pb"
	"testing"
	"time"
)

// 判断功能是否正确,channel是否能正常收发
func TestAppendEntriesRpc(t *testing.T) {
	nodes := RaftCreateCluster(t)
	if len(nodes) == 0 {
		log.Fatal("node num must larger than 0.")
	}
	node := nodes[0]
	go node.run()
	logDebug("TestAppendEntriesRpc(), raft info, append entriy addr:%p", node.rpc.rpcCh.rpcAppendEntryRequestCh)
	isPrintRpcAppendEntryRequestCh[0] = false
	datas := [4]string{"I am command1.", "I am command2.", "I am command3.", "I am command4."}
	req := &pb.AppendEntryRequest{
		Ver:          &RPOTO_VER_APPEND_ENTRY_REQUEST,
		LeaderTerm:   1,
		LeaderID:     "0",
		PrevLogIndex: MAX_LOG_INDEX_NUM,
		PrevLogTerm:  0,
		LogType:      uint32(LogCommand),
		LeaderCommit: MAX_LOG_INDEX_NUM,
	}
	// 发送四个请求，不走RPC
	for _, data := range datas {
		entryEncode, err := json.Marshal(createLogEntry(0, LogCommand, 1, []byte(data)))
		if err != nil {
			log.Fatal(err)
		}
		req.Entry = entryEncode
		node.rpc.rpcCh.rpcAppendEntryRequestCh <- req
		resp := <-node.rpc.rpcCh.rpcAppendEntryResponseCh
		logDebug("TestAppendEntriesRpc(), raft info, append entriy resp: %s", resp)
	}

}

func TestExecCommand(t *testing.T) {
	// 启动集群
	nodes := RaftCreateCluster(t)
	if len(nodes) == 0 {
		t.Fatalf("nodes must greater than 0.")
	}
	// 定义callbackFunc
	cbFunc := func(data []byte) error {
		log.Println("cbFunc:", string(data))
		return nil
	}
	for _, node := range nodes {
		node.storage.registerCallBackFunc(cbFunc)
	}
	// 定义datas
	datas := [4]string{"I am command1.", "I am command2.", "I am command3.", "I am command4."}
	nodes[0].storage.registerCallBackFunc(cbFunc)
	nodes[0].raftState.state = Leader
	nodes[0].raftState.currentTerm = 10
	nodes[0].setLeader(Servers[0])
	for _, node := range nodes {
		go node.run()
	}
	// 构造请求结构体
	req := &pb.ExecCommandRequest{
		Ver:     &PROTO_VER_EXEC_COMMAND_REQUEST,
		LogType: uint32(LogCommand),
	}
	for _, data := range datas {
		req.Command = []byte(data)
		nodes[0].rpc.rpcCh.rpcExecCommandRequestCh <- req
		logInfo("req to rpcExecCommandRequestCh:%v", req)

		resp := <-nodes[0].rpc.rpcCh.rpcExecCommandResponseCh
		logInfo("resp from rpcExecCommandResponseCh:%v", resp)
	}
	<-time.After(time.Second * 5)

}
