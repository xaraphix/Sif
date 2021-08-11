package raftadapter

import (
	. "github.com/xaraphix/Sif/internal/raft"
)

type RaftNodeAdapter struct {
	//TODO
}

func (a RaftNodeAdapter) RequestVoteFromPeer(peer Peer, voteRequest VoteRequest) VoteResponse {
	return VoteResponse{}
}

func (a RaftNodeAdapter) ReplicateLog(peer Peer, logRequest LogRequest) {
 return
}

func (a RaftNodeAdapter) ReceiveLogRequest(logRequest LogRequest) LogResponse {
 return LogResponse{ }
 }

func (a RaftNodeAdapter) GenerateVoteResponse(VoteRequest) VoteResponse {
	return VoteResponse{}
}
