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
