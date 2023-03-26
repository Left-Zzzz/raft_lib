package raftlib

import (
	"log"
	"testing"
	"time"
)

// 集群启动测试
func TestRaftCreateCluster(t *testing.T) {
	nodes := RaftCreateCluster(t)
	// 测试raft节点初始状态
	for _, node := range nodes {
		// leaderID初始值必须为空字符串
		if node.Leader() != "" {
			t.Fatalf("node's initial leader_id should be \"\"(empty string)\n")
		}
		// raft状态初始值必须为Follower
		if node.raftState.state != Follower {
			t.Fatalf("node's initial raftstate should be Follower\n")
		}
		//node.rpc.server.Stop()
	}
}

// 选举测试
func TestRaftElection(t *testing.T) {
	nodes := RaftCreateCluster(t)
	if len(nodes) == 0 {
		t.Fatalf("nodes must greater than 0.")
	}
	nodes[0].raftState.state = Candidate
	for _, node := range nodes {
		go node.run()
	}
	<-time.After(ElectionTimeout)
	for key, node := range nodes {
		log.Printf("node%d:state = %s, leader_id = %s.\n", key, node.getState(), string(node.Leader()))
	}

}

func RaftCreateCluster(t *testing.T) []*Raft {
	// TODO: TODO
	nodes := []*Raft{}
	for _, server := range Servers {
		node := createRaft(server, uint64(len(Servers)))
		nodes = append(nodes, node)
	}
	logInfo("Create cluster over.")
	return nodes
}

// 测试raft节点初始状态
func TestRaftNodeInfo(t *testing.T) {

}

// 测试有节点当选leader后其余节点状态是否节点
func TestNodeInfoChangeAfterElection(t *testing.T) {

	//node := runRaft(server, uint64(len(Servers)))
}
