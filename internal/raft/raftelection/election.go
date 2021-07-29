package raftelection

import (
	"github.com/xaraphix/Sif/internal/raft"
)

var (
	ElctnMgr raft.RaftElection = &ElectionManager{}
)

type ElectionManager struct {
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	rn.Mu.Lock()
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	rn.VotesReceived = nil
	rn.ElectionMonitor.Start(rn)
	rn.ElectionManager.RequestVotes(rn)
	rn.ElectionInProgress = false
	rn.Mu.Unlock()
}

func requestVoteFromPeer(rn *raft.RaftNode, vr raft.VoteRequest, voteResponseChannel chan raft.VoteResponse) {
	for _, peer := range rn.Peers {
		go func(p raft.Peer, vrc chan raft.VoteResponse) {
			vrc <- rn.RPCAdapter.RequestVoteFromPeer(p, vr)
		}(peer, voteResponseChannel)
	}

}
func (em *ElectionManager) RequestVotes(rn *raft.RaftNode) {

	voteResponseChannel := make(chan raft.VoteResponse)
	voteRequest := em.GenerateVoteRequest(rn)
	requestVoteFromPeer(rn, voteRequest, voteResponseChannel)
	concludeElection(rn, voteResponseChannel)

}

func concludeElection(rn *raft.RaftNode, voteResponseChannel chan raft.VoteResponse) {
	counter := 1

	for response := range voteResponseChannel {
		updateVotesRecieved(rn, response)
		if counter == len(rn.Peers) {
			break
		}
		counter = counter + 1
	}

	close(voteResponseChannel)

}

func (em *ElectionManager) StopElection(rn *raft.RaftNode) {

}

func (em *ElectionManager) GenerateVoteRequest(rn *raft.RaftNode) raft.VoteRequest {
	return raft.VoteRequest{}
}

func updateVotesRecieved(rn *raft.RaftNode, v raft.VoteResponse) {
	rn.VotesReceived = append(rn.VotesReceived, v.PeerId)
	//TODO
}
