package raftadapter

import "github.com/xaraphix/Sif/internal/raft"

type RaftNodeAdapter struct {
	//TODO
}

func (ra *RaftNodeAdapter) RequestVoteFromPeer(peer int32, voteRequest raft.VoteRequest) raft.VoteResponse {
	//TODO
	return raft.VoteResponse{}
}

func RequestVoteFromPeer(peer *raft.Peer, voteRequest raft.VoteRequest) raft.VoteResponse {
	return raft.VoteResponse{PeerId: 2, VoteGranted: false}
}
