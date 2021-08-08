package raftelection

import (
	"math/rand"
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)

var (
	ElctnMgr raft.RaftElection = &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		ElectionTimerOff: true,
	}
)

type ElectionManager struct {
	ElectionTimeoutDuration time.Duration
	ElectionTimerOff bool
}

func (el *ElectionManager) HasElectionTimerStarted() bool {
	return !el.ElectionTimerOff
}

func (el *ElectionManager) HasElectionTimerStopped() bool {
	return el.ElectionTimerOff
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	electionChannel := make(chan raft.ElectionUpdates)
	done := make(chan bool)
	defer close(done)
	
	rn.Mu.Lock()
	defer rn.Mu.Unlock()

	rn.LeaderHeartbeatMonitor.Stop()
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	rn.VotesReceived = nil

	for {
		select {
		case <- followerAccordingToPeer:
			becomeAFollowerAccordingToPeer()
		case <- leader:
			becomeALeader()
		case <- leaderHeartbeatReceived:
			becomeAFollowerAccordingToLeader()
		case <- electionTimedOut:
			restartElection()
		}
	}
}

func (em *ElectionManager) RequestVotes(
	rn *raft.RaftNode,
	electionChannel chan raft.ElectionUpdates) {

	voteResponseChannel := make(chan raft.VoteResponse)
	voteRequest := em.GenerateVoteRequest(rn)
	requestVoteFromPeer(rn, voteRequest, voteResponseChannel)
	go em.setupElectionTimer(rn, electionChannel)
	startTimer(electionChannel)
	concludeElection(rn, voteResponseChannel, electionChannel)
}

func startTimer(ec chan raft.ElectionUpdates) {
	eu := <- ec
	if !eu.ElectionTimerStarted {
		panic("Election Timer not started") 
	}
}


func (em *ElectionManager) GetResponseForVoteRequest(rn *raft.RaftNode, vr raft.VoteRequest) raft.VoteResponse {
 return	getVoteResponseForVoteRequest(rn, vr)
}

func (em *ElectionManager) setupElectionTimer(
	rn *raft.RaftNode,
	electionChannel chan raft.ElectionUpdates) {
	eu := raft.ElectionUpdates{
		ElectionTimerStarted: true,
	}

	electionChannel <- eu
	em.ElectionTimerOff = false

	eu = raft.ElectionUpdates{
		ElectionStopped: false,
	}

	time.Sleep(em.ElectionTimeoutDuration)
	if rn.ElectionInProgress == true && eu.ElectionStopped == false {
		//kill the Current Election
		electionChannel <- raft.ElectionUpdates{ElectionOvertimed: true}
		rn.ElectionInProgress = false
		rn.CurrentRole = raft.FOLLOWER
		rn.VotesReceived = nil
		go em.StartElection(rn)
	}

	em.ElectionTimerOff = true
	close(electionChannel)
}

func requestVoteFromPeer(
	rn *raft.RaftNode,
	vr raft.VoteRequest,
	voteResponseChannel chan raft.VoteResponse) {

	for _, peer := range rn.Peers {
		go func(p raft.Peer, vrc chan raft.VoteResponse) {
			
			voteResponse := rn.RPCAdapter.RequestVoteFromPeer(p, vr)
			vrc <- voteResponse
		}(peer, voteResponseChannel)
	}
}

