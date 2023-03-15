package raftlib

import (
	"log"
	"raft_lib/pb"
	"testing"
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
	req := &pb.AppendEntryRequest{}
	node.rpc.rpcCh.rpcAppendEntryRequestCh <- req
	resp := <-node.rpc.rpcCh.rpcAppendEntryResponseCh
	logDebug("TestAppendEntriesRpc(), raft info, append entriy resp: %s", resp)
}

func TestExecCommand(t *testing.T) {
	// 启动集群
	nodes := RaftCreateCluster(t)
	if len(nodes) == 0 {
		t.Fatalf("nodes must greater than 0.")
	}
	nodes[0].raftState.state = Leader
	nodes[0].raftState.currentTerm = 10
	nodes[0].setLeader(Servers[0])
	for _, node := range nodes {
		go node.run()
	}
	// 构造请求结构体
	req := &pb.ExecCommandRequest{
		Ver:     &PROTO_VER_EXEC_COMMAND_REQUEST,
		Command: []byte("I am command."),
	}
	nodes[0].rpc.rpcCh.rpcExecCommandRequestCh <- req
	logInfo("req to rpcExecCommandRequestCh:%v", req)

	resp := <-nodes[0].rpc.rpcCh.rpcExecCommandResponseCh
	logInfo("resp from rpcExecCommandResponseCh:%v", resp)

}
