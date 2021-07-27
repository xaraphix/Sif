package raft

type VoteRequest struct {
	//TODO
}

type VoteResponse struct {
	//TODO
}

func startElection(rn *RaftNode) {
	rn.Mu.Lock()
	rn.ElectionInProgress = true
	rn.CurrentRole = CANDIDATE
	rn.VotedFor = raftnode.Id
	rn.CurrentTerm = raftnode.CurrentTerm + 1
	rn.VotesReceived = nil
	rn.StartElectionMonitor()
	rn.RequestVotes()
	rn.ElectionInProgress = false
	rn.Mu.Unlock()

}
