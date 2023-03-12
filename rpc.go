package raftlib

import (
	"context"
	"errors"
	"log"
	"net"
	"raft_lib/pb"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// RPC层和Raft层通过通道将RPC消息传输，RPC层只负责消息收发，Raft层负责处理消息
type RpcCh struct {
	// RequestVoteRPC channel
	rpcRequestVoteRequestCh  chan *pb.RequestVoteRequest
	rpcRequestVoteResponseCh chan *pb.RequestVoteResponse
	// AppendEntryRPC channel
	rpcAppendEntryRequestCh  chan *pb.AppendEntryRequest
	rpcAppendEntryResponseCh chan *pb.AppendEntryResponse
}

type Rpc struct {
	sync.RWMutex
	server      *grpc.Server
	clientConns map[ServerID]*grpc.ClientConn
	clients     map[ServerID]pb.RpcServiceClient
	rpcCh       *RpcCh

	// 继承 protoc-gen-go-grpc 生成的服务端代码
	pb.UnimplementedRpcServiceServer
}

func (r *Rpc) createRpcServer(raft *Raft) (*grpc.Server, error) {
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	// 实例化RpcCh对象
	rpcCh := &RpcCh{
		rpcRequestVoteRequestCh:  make(chan *pb.RequestVoteRequest, 1),
		rpcRequestVoteResponseCh: make(chan *pb.RequestVoteResponse, 1),
		rpcAppendEntryRequestCh:  make(chan *pb.AppendEntryRequest, 1),
		rpcAppendEntryResponseCh: make(chan *pb.AppendEntryResponse, 1),
	}
	logDebug("createRpcServer():rpcCh:%v", *rpcCh)
	// 将实现的接口注册进 gRPC 服务器
	pb.RegisterRpcServiceServer(s, &Rpc{
		rpcCh: rpcCh,
	})
	r.rpcCh = rpcCh
	return s, nil
}

// 启动rpc服务器
func (r *Rpc) runRpcServer(s *grpc.Server, port ServerPort) {
	// 监听本地 port 端口
	port_formed := string(":" + port)
	listen, err := net.Listen("tcp", port_formed)
	if err != nil {
		panic(err)
	}
	logInfo("gRPC server starts running...")
	// 启动 gRPC 服务器
	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}

// 启动rpc客户端
func (r *Rpc) runRpcClient(peer Server) error {
	// 如果已经创建过了，直接返回
	if clientConnTemp, ok := r.clientConns[peer.ID]; ok {
		// 如果断线，尝试重连
		if clientConnTemp.GetState() != connectivity.Ready {
			clientConnTemp.Connect()
		}
		// 如果还断线，报错
		if clientConnTemp.GetState() != connectivity.Ready {
			return errors.New("grpc client can not connect to server")
		}
		logDebug("gRPC client have been created.\n")
		return nil
	}
	logDebug("Creating gRPC client.\n")
	// 创建客户端
	clientConn, err := grpc.Dial(string(peer.Address)+":"+string(peer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 断线重连参数
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true}))
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 创建RpcService
	clientService := pb.NewRpcServiceClient(clientConn)
	r.clientConns[peer.ID] = clientConn
	r.clients[peer.ID] = clientService
	return nil
}

// 处理RequestVoteRpc请求并返回结果
func (r *Rpc) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	logDebug("start process requestvote.\n")
	// 将消息传送给channel
	logDebug("RequestVote(): request vote request channel addr:%p", r.rpcCh.rpcRequestVoteRequestCh)
	r.rpcCh.rpcRequestVoteRequestCh <- req
	resp := <-r.rpcCh.rpcRequestVoteResponseCh
	logDebug("rpcRequestVoteResponseCh resp:%v", resp)
	return resp, nil
}

// 处理AppendEntryRpc请求并返回结果
func (r *Rpc) AppendEntry(ctx context.Context, req *pb.AppendEntryRequest) (*pb.AppendEntryResponse, error) {
	logDebug("start process append entry.\n")
	// 将消息传送给channel
	logDebug("AppendEntry(): append entry request channel addr:%p", r.rpcCh.rpcAppendEntryRequestCh)
	r.rpcCh.rpcAppendEntryRequestCh <- req
	resp := <-r.rpcCh.rpcAppendEntryResponseCh
	logDebug("rpcAppendEntryResponseCh resp:%v", resp)
	return resp, nil
}
