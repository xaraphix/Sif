package raft

type Peer struct {
	Id      int32
	Address string
}

func (p *Peer) RequestVote(peerOf *RaftNode, voteRequest VoteRequest) {
	voteResponse := raftAdapter.RequestVoteFromPeer(p, voteRequest)
	peerOf.UpdateVotesRecieved(voteResponse)
}

func requestVotesFromPeers(rn *RaftNode) {
	voteRequest := rn.GenerateVoteRequest()
	for _, peer := range rn.Peers {
		go peer.RequestVote(rn, voteRequest)
	}

}
