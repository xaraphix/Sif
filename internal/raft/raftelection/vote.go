package raftelection

import "github.com/xaraphix/Sif/internal/raft"

func (em *ElectionManager) GenerateVoteRequest(rn *raft.RaftNode) raft.VoteRequest {
	return raft.VoteRequest{
		NodeId:      rn.Id,
		CurrentTerm: rn.CurrentTerm,
		LogLength:   int32(len(rn.Logs)),
		LastTerm:    rn.LogMgr.GetLog(rn, int32(len(rn.Logs)-1)).Term,
	}
}

func becomeAFollowerAccordingToPeersTerm(
	rn *raft.RaftNode,
	v raft.VoteResponse,
	electionUpdates *raft.ElectionUpdates) {

	rn.CurrentTerm = v.Term
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.ElectionInProgress = false
	rn.LeaderHeartbeatMonitor.Start(rn)
}

func (em *ElectionManager) getVoteResponseForVoteRequest(rn *raft.RaftNode, voteRequest raft.VoteRequest) raft.VoteResponse {
	voteResponse := raft.VoteResponse{}
	logOk := isCandidateLogOK(rn, voteRequest)
	termOK := isCandidateTermOK(rn, voteRequest)

	if logOk && termOK {
		if rn.ElectionInProgress {
			em.peerVoteChannel <- raft.RaftNode{
				Node: raft.Node{
					Id:          voteRequest.NodeId,
					CurrentTerm: voteRequest.CurrentTerm,
				},
			}
		}
		rn.CurrentTerm = voteRequest.CurrentTerm
		rn.CurrentRole = raft.FOLLOWER
		rn.VotedFor = voteRequest.NodeId
		voteResponse.PeerId = rn.Id
		voteResponse.Term = rn.CurrentTerm
		voteResponse.VoteGranted = true
	} else {
		voteResponse.PeerId = rn.Id
		voteResponse.Term = rn.CurrentTerm
		voteResponse.VoteGranted = false
	}

	return voteResponse
}

func isCandidateLogOK(rn *raft.RaftNode, vr raft.VoteRequest) bool {

	myLogTerm := rn.LogMgr.GetLog(rn, int32(len(rn.Logs)-1)).Term
	logOk := vr.LastTerm > myLogTerm ||
		(vr.LastTerm == myLogTerm && vr.LogLength >= int32(len(rn.Logs)))
	return logOk
}

func isCandidateTermOK(rn *raft.RaftNode, vr raft.VoteRequest) bool {
	termOk := vr.CurrentTerm > rn.CurrentTerm ||
		(vr.CurrentTerm == rn.CurrentTerm && hasCandidateBeenVotedPreviously(rn, vr))

	return termOk
}

func hasCandidateBeenVotedPreviously(rn *raft.RaftNode, voteRequest raft.VoteRequest) bool {
	return false
}

func (em *ElectionManager) concludeFromReceivedVotes(rn *raft.RaftNode) {

	em.followerAccordingToPeer = make(chan raft.VoteResponse)
	em.leader = make(chan bool)
	em.follower = make(chan bool)

	go func() {
		for {
			if em.VotesReceived == nil {
				em.VotesReceived = []raft.VoteResponse{}
			}
			if len(em.VotesReceived) == len(rn.Peers) {
				continue
			}
			for vr := range em.votesResponse {
				em.VotesReceived = append(em.VotesReceived, vr)
				em.concludeFromReceivedVote(rn, vr)
				if len(em.VotesReceived) == len(rn.Peers) {
					break
				}
			}
		}
	}()
}

func (em *ElectionManager) concludeFromReceivedVote(rn *raft.RaftNode, vr raft.VoteResponse) {
	valid := isElectionValid(rn, vr)
	if !valid && isPeerTermHigher(rn, vr) {
		em.followerAccordingToPeer <- vr
	} else if becomeLeader, becomeFollower := em.whatDoesTheMajorityWant(len(rn.Peers), vr); true {
		if becomeLeader {
			em.leader <- true
		} else if becomeFollower {
			em.follower <- true
		}
	}

}

func isPeerTermHigher(rn *raft.RaftNode, vr raft.VoteResponse) bool {
	return rn.CurrentTerm < vr.Term
}

func isElectionValid(rn *raft.RaftNode, vr raft.VoteResponse) bool {
	if rn.CurrentRole == raft.CANDIDATE &&
		rn.CurrentTerm == vr.Term {
		return true
	} else {
		return false
	}
}

func (em *ElectionManager) whatDoesTheMajorityWant(numOfPeers int, vr raft.VoteResponse) (bool, bool) {

	majorityCount := numOfPeers / 2
	votesInFavor := 0
	leader, follower := false, false
	for _, vr := range em.VotesReceived {
		if vr.VoteGranted {
			votesInFavor = votesInFavor + 1
		}
	}

	// if majority has voted against, game over
	if len(em.VotesReceived)-votesInFavor > majorityCount {
		follower = true
	}

	// if majority has voted in favor, game won
	if votesInFavor >= majorityCount {
		leader = true
	}

	return leader, follower
}
