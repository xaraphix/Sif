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
		concludeElectionIfPossible(rn, response)
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

func concludeElectionIfPossible(rn *raft.RaftNode, v raft.VoteResponse) {
	rn.VotesReceived = append(rn.VotesReceived, v)
	//TODO
	majorityCount := len(rn.Peers) / 2
	votesInFavor := 0

	for _, vr := range rn.VotesReceived {
		if vr.VoteGranted {
			votesInFavor = votesInFavor + 1
		}
	}

	// if majority has voted against, game over
	if len(rn.VotesReceived)-votesInFavor >= majorityCount {
		rn.Mu.Unlock()
		rn.CurrentRole = raft.FOLLOWER
		rn.Mu.Lock()
	}

	// if majority has voted in favor, game won
	if votesInFavor >= majorityCount {
		rn.Mu.Unlock()
		rn.CurrentRole = raft.LEADER
		rn.Mu.Lock()
	}
}
