package raftadapter

import (
	. "github.com/xaraphix/Sif/internal/raft"
	"github.com/xaraphix/Sif/internal/raft/protos"
)

type RaftNodeAdapter struct {
	//TODO
}

func (a RaftNodeAdapter) RequestVoteFromPeer(peer Peer, voteRequest *protos.VoteRequest) *protos.VoteResponse {
	return &protos.VoteResponse{}
}

func (a RaftNodeAdapter) ReplicateLog(peer Peer, logRequest *protos.LogRequest) {
 return
}

func (a RaftNodeAdapter) ReceiveLogRequest(logRequest *protos.LogRequest) *protos.LogResponse {
 return &protos.LogResponse{ }
 }

func (a RaftNodeAdapter) GenerateVoteResponse(*protos.VoteRequest) *protos.VoteResponse {
	return &protos.VoteResponse{}
}
