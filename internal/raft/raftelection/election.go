package raftelection

import (
	"math/rand"
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)

var (
	ElctnMgr raft.RaftElection = &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
	}
)

type ElectionManager struct {
	ElectionTimeoutDuration time.Duration
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	electionChannel := make(chan raft.ElectionUpdates)

	rn.Mu.Lock()
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	rn.VotesReceived = nil

	go em.restartElectionWhenItTimesOut(rn, electionChannel)
	rn.ElectionMgr.RequestVotes(rn, electionChannel)
	rn.ElectionInProgress = false

	ifLeaderStartHeartbeatTransmitter(rn)

	rn.Mu.Unlock()
}

func (em *ElectionManager) restartElectionWhenItTimesOut(rn *raft.RaftNode, electionChannel chan raft.ElectionUpdates) {
	time.Sleep(em.ElectionTimeoutDuration)
	if rn.ElectionInProgress == true {
		//kill the Current Election
		electionChannel <- raft.ElectionUpdates{ElectionOvertimed: true}
		rn.Mu.Unlock()
		rn.ElectionInProgress = false
		rn.CurrentRole = raft.FOLLOWER
		rn.VotesReceived = nil
		go em.StartElection(rn)
	}
	close(electionChannel)
}

func (em *ElectionManager) RequestVotes(rn *raft.RaftNode, electionChannel chan raft.ElectionUpdates) {
	voteResponseChannel := make(chan raft.VoteResponse)
	voteRequest := em.GenerateVoteRequest(rn)
	requestVoteFromPeer(rn, voteRequest, voteResponseChannel)
	concludeElection(rn, voteResponseChannel, electionChannel)
}

func requestVoteFromPeer(rn *raft.RaftNode, vr raft.VoteRequest, voteResponseChannel chan raft.VoteResponse) {
	for _, peer := range rn.Peers {
		go func(p raft.Peer, vrc chan raft.VoteResponse) {
			vrc <- rn.RPCAdapter.RequestVoteFromPeer(p, vr)
		}(peer, voteResponseChannel)
	}
}

func concludeElection(rn *raft.RaftNode, voteResponseChannel chan raft.VoteResponse, electionChannel chan raft.ElectionUpdates) {
	counter := 1
	electionUpdates := &raft.ElectionUpdates{ElectionOvertimed: false}

	go func(ec chan raft.ElectionUpdates) {
		updates := <-ec
		electionUpdates.ElectionOvertimed = updates.ElectionOvertimed
	}(electionChannel)

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

}

func (em *ElectionManager) GenerateVoteRequest(rn *raft.RaftNode) raft.VoteRequest {
	return raft.VoteRequest{}
}


func concludeElectionIfPossible(rn *raft.RaftNode, v raft.VoteResponse, electionUpdates *raft.ElectionUpdates) {
	if voteResponseConsidersElectionValid(rn, v) {
		concludeFromVoteIfPossible(rn, v, electionUpdates)
	} else {
		becomeAFollowerAccordingToPeersTerm(rn, v)
	}
}


func voteResponseConsidersElectionValid(rn *raft.RaftNode, v raft.VoteResponse) bool {
	if rn.CurrentRole == raft.CANDIDATE &&
	rn.CurrentTerm == v.CurrentTerm &&
	v.VoteGranted {
		return true
	} else {
		return false
	}
}


func concludeFromVoteIfPossible(rn *raft.RaftNode, v raft.VoteResponse, electionUpdates *raft.ElectionUpdates){
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
		rn.Mu.Lock()
}

func becomeTheLeader(rn *raft.RaftNode) {
		rn.Mu.Unlock()
		rn.CurrentRole = raft.LEADER
		rn.Mu.Lock()
}

func becomeAFollowerAccordingToPeersTerm(rn *raft.RaftNode, v raft.VoteResponse) {
	rn.Mu.Unlock()
	rn.CurrentTerm = v.CurrentTerm
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.Mu.Lock()
	//cancel the election
}
