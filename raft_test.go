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
	nodes[0].setState(Candidate)
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

// NoOp测试
func TestRaftNoOp(t *testing.T) {
	nodes := RaftCreateCluster(t)
	if len(nodes) == 0 {
		t.Fatalf("nodes must greater than 0.")
	}
	// 手动设置节点0为leader
	nodes[0].setState(Leader)
	nodes[0].setLeader(Servers[0])
	nodes[0].setCurrentTerm(5)
	// 判断NoOp补丁是否有效, 即是否能将之前未提交的日志安全提交

	// 构造5个已提交的日志，并追加日志后提交
	for i := 0; i < 5; i++ {
		entry := Log{
			Index:   uint32(i),
			Term:    1,
			LogType: LogCommand,
			Data:    []byte("i am a data."),
		}
		nodes[0].storage.appendEntryEntity(uint32(i), entry)
		nodes[0].storage.commit(uint32(i), func([]byte) error { return nil })
	}
	nodes[0].setCommitIndex(4)
	nodes[0].setLastEntry(4, 1)

	// 定义callbackFunc
	cbFunc := func(data []byte) error {
		log.Println("cbFunc:", string(data))
		return nil
	}
	for _, node := range nodes {
		node.storage.registerCallBackFunc(cbFunc)
	}
	// 运行节点
	for _, node := range nodes {
		go node.run()
	}
	// 等待两倍心跳超时时间
	<-time.After(HeartbeatTimeout * 2)
	// 获取每个节点下标为0的日志
	entrys := []Log{}
	for _, node := range nodes {
		entry, err := node.storage.getEntry(5)
		if err != nil {
			t.Fatal(err)
		}
		entrys = append(entrys, entry)
	}
	// 判断NoOp补丁是否成功加入到日志
	for _, entry := range entrys {
		if entry.LogType != LogNoOp {
			t.Fail()
		}
	}

}
