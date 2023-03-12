package raftlib

import (
	"context"
	"raft_lib/pb"
)

// RequestVote RPC请求流程
func (r *Rpc) RequestVoteRequest(
	peer Server, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	// 尝试启动rpc
	err := r.runRpcClient(peer)
	if err != nil {
		return nil, err
	}
	logDebug("start send requestvote.\n")

	ctx := context.Background()
	var resp *pb.RequestVoteResponse
	resp, err = r.clients[peer.ID].RequestVote(ctx, req)
	if err != nil {
		logError(": An error occured while calling RequestVoteRPC: ", err)
		return nil, err
	}

	return resp, err
}

// AppendEntry RPC请求流程
func (r *Rpc) AppendEntryRequest(
	peer Server, req *pb.AppendEntryRequest) (*pb.AppendEntryResponse, error) {
	// 尝试启动rpc
	err := r.runRpcClient(peer)
	if err != nil {
		return nil, err
	}
	logDebug("start send append entry request.\n")

	ctx := context.Background()
	var resp *pb.AppendEntryResponse
	resp, err = r.clients[peer.ID].AppendEntry(ctx, req)
	if err != nil {
		logError(": An error occured while calling AppendEntryRPC: ", err)
		return nil, err
	}

	return resp, err
}
