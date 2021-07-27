package raft

type Peer struct {
	Id      int32
	Address string
}

func requestVotesFromPeers(rn *RaftNode, voteReponseChannel chan VoteResponse) {
	voteRequest := rn.GenerateVoteRequest()

	for _, peer := range rn.Peers {
		go func(p *Peer, vrc chan VoteResponse) {
			vrc <- raftAdapter.RequestVoteFromPeer(p, voteRequest)
		}(&peer, voteReponseChannel)
	}
}
