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
	rn.ElectionManager.RequestVotes(rn, electionChannel)
	rn.ElectionInProgress = false

	ifLeaderStartHeartbeatTransmitter(rn)

	rn.Mu.Unlock()
}

func ifLeaderStartHeartbeatTransmitter(rn *raft.RaftNode) {

}

func (em *ElectionManager) restartElectionWhenItTimesOut(rn *raft.RaftNode, electionChannel chan raft.ElectionUpdates) {
	time.Sleep(em.ElectionTimeoutDuration)
	if rn.ElectionInProgress == true {
		//kill the Current Election
		electionChannel <- raft.ElectionUpdates{ElectionOvertimed: true}
		em.StartElection(rn)
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
			rn.VotesReceived = nil
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
	rn.VotesReceived = append(rn.VotesReceived, v)
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
		rn.Mu.Unlock()
		rn.CurrentRole = raft.FOLLOWER
		rn.Mu.Lock()
	}

	// if majority has voted in favor, game won
	if votesInFavor >= majorityCount &&
		electionUpdates.ElectionOvertimed == false {
		rn.Mu.Unlock()
		rn.CurrentRole = raft.LEADER
		rn.Mu.Lock()
	}
}
