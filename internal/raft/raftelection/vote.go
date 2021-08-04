package raftelection

import "github.com/xaraphix/Sif/internal/raft"

func concludeElection(
	rn *raft.RaftNode,
	voteResponseChannel chan raft.VoteResponse,
	electionOvertimeChannel <-chan raft.ElectionUpdates) {

	counter := 1
	electionUpdates := &raft.ElectionUpdates{ElectionOvertimed: false}

	go func(ec <-chan raft.ElectionUpdates) {
		updates := <-ec
		electionUpdates.ElectionOvertimed = updates.ElectionOvertimed
	}(electionOvertimeChannel)

	for response := range voteResponseChannel {
		if electionUpdates.ElectionOvertimed == false {
			concludeElectionIfPossible(rn, response, electionUpdates)
			if counter == len(rn.Peers) {
				break
			}
			counter = counter + 1
		} else {
			rn.Mu.Unlock()
			rn.VotesReceived = nil
			rn.Mu.Lock()
		}
	}

	close(voteResponseChannel)

}

func (em *ElectionManager) StopElection(rn *raft.RaftNode) {
	rn.ElectionInProgress = false
}

func (em *ElectionManager) GenerateVoteRequest(rn *raft.RaftNode) raft.VoteRequest {
	return raft.VoteRequest{
		NodeId:      rn.Id,
		CurrentTerm: rn.CurrentTerm,
		LogLength:   int32(len(rn.Logs)),
		LastTerm:    rn.LogMgr.GetLog(int32(len(rn.Logs) - 1)).Term,
	}
}

func concludeElectionIfPossible(
	rn *raft.RaftNode,
	v raft.VoteResponse,
	electionUpdates *raft.ElectionUpdates) {

	if voteResponseConsidersElectionValid(rn, v) {
		concludeFromVoteIfPossible(rn, v, electionUpdates)
	} else if v.Term > rn.CurrentTerm {
		becomeAFollowerAccordingToPeersTerm(rn, v, electionUpdates)
	}
}

func voteResponseConsidersElectionValid(rn *raft.RaftNode, v raft.VoteResponse) bool {
	if rn.CurrentRole == raft.CANDIDATE &&
		rn.CurrentTerm == v.Term &&
		v.VoteGranted {
		return true
	} else {
		return false
	}
}

func concludeFromVoteIfPossible(
	rn *raft.RaftNode,
	v raft.VoteResponse,
	electionUpdates *raft.ElectionUpdates) {

	majorityCount := len(rn.Peers) / 2
	votesInFavor := 0

	for _, vr := range rn.VotesReceived {
		if vr.VoteGranted {
			votesInFavor = votesInFavor + 1
		}
	}

	// if majority has voted against, game over
	if len(rn.VotesReceived)-votesInFavor > majorityCount &&
		electionUpdates.ElectionOvertimed == false {
		becomeAFollower(rn)
	}

	// if majority has voted in favor, game won
	if votesInFavor >= majorityCount &&
		electionUpdates.ElectionOvertimed == false {
		becomeTheLeader(rn)
	}
}

func becomeAFollower(rn *raft.RaftNode) {
	rn.Mu.Unlock()
	rn.CurrentRole = raft.FOLLOWER
	rn.ElectionInProgress = false
	rn.Mu.Lock()
}

func becomeTheLeader(rn *raft.RaftNode) {
	rn.Mu.Unlock()
	rn.CurrentRole = raft.LEADER
	rn.ElectionInProgress = false
	rn.Mu.Lock()
}

func becomeAFollowerAccordingToPeersTerm(
	rn *raft.RaftNode,
	v raft.VoteResponse,
	electionUpdates *raft.ElectionUpdates) {

	rn.Mu.Unlock()
	rn.CurrentTerm = v.Term
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.ElectionInProgress = false
	rn.Mu.Lock()
}

func getVoteResponseForVoteRequest(rn *raft.RaftNode, voteRequest raft.VoteRequest) raft.VoteResponse {
	voteResponse := raft.VoteResponse{}
	logOk := isCandidateLogOK(rn, voteRequest)
	termOK := isCandidateTermOK(rn, voteRequest)

	if logOk && termOK {
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

	myLogTerm := rn.LogMgr.GetLog(int32(len(rn.Logs) - 1)).Term
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
