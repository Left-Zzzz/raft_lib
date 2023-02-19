package raftlib

import (
	"context"
	"log"
	"net"
	"raft_lib/pb"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Rpc struct {
	sync.RWMutex
	server *grpc.Server
	client *grpc.ClientConn

	// 继承 protoc-gen-go-grpc 生成的服务端代码
	pb.UnimplementedRpcServiceServer
}

// 服务rpc通用请求响应流程
func (rpc *Rpc) rpcCall(ctx context.Context, in *pb.ReqInfoBase) (*pb.RspInfoBase, error) {
	log.Println("client call rpcCall...")
	// TODO：通用处理逻辑
	log.Println(in)
	return &pb.RspInfoBase{}, nil
}

// 启动rpc服务器
func (r *Rpc) runRpcServer(addr ServerAddress, port ServerPort) (*grpc.Server, error) {
	// 监听本地 port 端口
	port_formed := string(":" + port)
	listen, err := net.Listen("tcp", port_formed)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterRpcServiceServer(s, &Rpc{})
	log.Println("gRPC server starts running...")
	// 启动 gRPC 服务器
	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return s, nil
}

// 启动rpc客户端
func (r *Rpc) runRpcClient(addr ServerAddress, port ServerPort) {
	client, err := grpc.Dial("127.0.0.1:5678", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	r.client = client
}

// RequestVote RPC流程
func (r *Rpc) RequestVote(id ServerID, addr ServerAddress, req *pb.VoteRequest, resp *pb.VoteResponse) error {
	return nil
}
