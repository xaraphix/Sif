package raftadapter

import (
	. "github.com/xaraphix/Sif/internal/raft"
	pb "github.com/xaraphix/Sif/internal/raft/protos"
)

type RaftNodeAdapter struct {
	//TODO
}

func (a RaftNodeAdapter) RequestVoteFromPeer(peer Peer, voteRequest *pb.VoteRequest) *pb.VoteResponse {
	return &pb.VoteResponse{}
}

func (a RaftNodeAdapter) ReplicateLog(peer Peer, logRequest *pb.LogRequest) {
 return
}

func (a RaftNodeAdapter) ReceiveLogRequest(logRequest *pb.LogRequest) *pb.LogResponse {
 return &pb.LogResponse{ }
 }

func (a RaftNodeAdapter) GenerateVoteResponse(*pb.VoteRequest) *pb.VoteResponse {
	return &pb.VoteResponse{}
}
