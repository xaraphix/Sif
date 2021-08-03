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

func (em *ElectionManager) restartElectionWhenItTimesOut(
	rn *raft.RaftNode,
	electionChannel chan<- raft.ElectionUpdates)	{
	eu := raft.ElectionUpdates{
		ElectionStopped: false,
	}

	time.Sleep(em.ElectionTimeoutDuration)
	if rn.ElectionInProgress == true && eu.ElectionStopped == false {
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

func (em *ElectionManager) RequestVotes(
	rn *raft.RaftNode,
	electionChannel <-chan raft.ElectionUpdates,){

	voteResponseChannel := make(chan raft.VoteResponse)
	voteRequest := em.GenerateVoteRequest(rn)
	requestVoteFromPeer(rn, voteRequest, voteResponseChannel)
	concludeElection(rn, voteResponseChannel, electionChannel)
}

func requestVoteFromPeer(
	rn *raft.RaftNode,
	vr raft.VoteRequest,
	voteResponseChannel chan raft.VoteResponse) {

	for _, peer := range rn.Peers {
		go func(p raft.Peer, vrc chan raft.VoteResponse) {
			vrc <- rn.RPCAdapter.RequestVoteFromPeer(p, vr)
		}(peer, voteResponseChannel)
	}
}
