package raftlib

import (
	"context"
	"raftlib/pb"
)

// RequestVote RPC请求流程
func (r *Rpc) RequestVoteRequest(
	peer Server, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	resp := &pb.RequestVoteResponse{
		VoteGranted: false,
	}
	// 尝试启动rpc
	err := r.runRpcClient(peer)
	if err != nil {
		return resp, err
	}
	logDebug("start send requestvote.\n")

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), r.config.RpcTimeout)
	defer cancel()

	r.clientsMutex.Lock()
	// 当rpc调用失败是rpcResp会返回空，后续对resp的操作会出错，这里使用rpcResp替换掉resp返回
	rpcResp, err := r.clients[peer.ID].RequestVote(ctx, req)
	r.clientsMutex.Unlock()
	if err != nil {
		logError("An error occured while calling RequestVoteRPC: ", err)
		return resp, err
	}

	return rpcResp, err
}

// AppendEntry RPC请求流程
func (r *Rpc) AppendEntryRequest(
	peer Server, req *pb.AppendEntryRequest) (*pb.AppendEntryResponse, error) {
	resp := &pb.AppendEntryResponse{
		Success: false,
	}
	// 尝试启动rpc
	err := r.runRpcClient(peer)
	if err != nil {
		return resp, err
	}

	logDebug("start send append entry request.\n")

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), r.config.RpcTimeout)
	defer cancel()

	r.clientsMutex.Lock()
	rpcResp, err := r.clients[peer.ID].AppendEntry(ctx, req)
	r.clientsMutex.Unlock()
	if err != nil {
		logError("An error occured while calling AppendEntryRPC:%v", err)
		return resp, err
	}

	return rpcResp, err
}
