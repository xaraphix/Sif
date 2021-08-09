package raftadapter

import "github.com/xaraphix/Sif/internal/raft"

type RaftNodeAdapter struct {
	//TODO
}

func (ra *RaftNodeAdapter) RequestVoteFromPeer(peer int32, voteRequest raft.VoteRequest) raft.VoteResponse {
	//TODO
	return raft.VoteResponse{}
}


func (ra *RaftNodeAdapter) GenerateVoteRequest(req raft.VoteRequest) raft.VoteResponse {
	return raft.Sif.ElectionMgr.GetResponseForVoteRequest(raft.Sif, req)
}

