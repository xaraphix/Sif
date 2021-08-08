package raftelection

import (
	"math/rand"
	"time"

	"github.com/xaraphix/Sif/internal/raft"
)

var (
	ElctnMgr raft.RaftElection = &ElectionManager{
		ElectionTimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		ElectionTimerOff:        true,
	}
)

type ElectionManager struct {
	ElectionTimeoutDuration time.Duration
	ElectionTimerOff        bool
}

func (el *ElectionManager) HasElectionTimerStarted() bool {
	return !el.ElectionTimerOff
}

func (el *ElectionManager) HasElectionTimerStopped() bool {
	return el.ElectionTimerOff
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	done := make(chan bool)
	electionTimerDone := make(chan bool)
	voteRequestsDone := make(chan bool)
	conclusionDone := make(chan bool)
	defer close(done)

	rn.Mu.Lock()
	defer rn.Mu.Unlock()

	rn.LeaderHeartbeatMonitor.Stop()
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	rn.VotesReceived = nil

	votesReceived := em.askForVotes(voteRequestsDone, rn)
	electionTimedOut := em.startElectionTimer(electionTimerDone, rn)
	followerAccordingToPeer, leader, follower := concludeFromReceivedVotes(conclusionDone, rn, votesReceived)

	for {
		select {
		case <-followerAccordingToPeer:
			becomeAFollowerAccordingToPeer()
		case <-leader:
			becomeALeader()
		// case <- leaderHeartbeatReceived:
		// becomeAFollowerAccordingToLeader()
		case <-follower:
			becomeAFollower()
		case <-electionTimedOut:
			restartElection()
		}
	}
}

func becomeAFollower() {

}

func concludeFromReceivedVotes(
	done chan bool,
	rn *raft.RaftNode,
	votesReceived chan raft.VoteResponse) (fatp chan raft.VoteResponse, l chan bool, f chan bool) {

	followerAccordingToPeer := make(chan raft.VoteResponse)
	leader := make(chan bool)
	follower := make(chan bool)

	go func() {
		defer close(followerAccordingToPeer)
		defer close(leader)
		defer close(follower)
		for {
			select {
			case <-done:
				break
			default:
				for voteResponse := range votesReceived {
					valid := isElectionValid(rn, voteResponse)
					if !valid {
						followerAccordingToPeer <- voteResponse
					} else if becomeLeader, becomeFollower := whatDoesTheMajorityWant(rn, voteResponse); true {
						if becomeLeader {
							leader <- true
						} else if becomeFollower {
							follower <- becomeFollower
						}
					}

				}
			}
		}
	}()

	return followerAccordingToPeer, leader, follower
}

func isElectionValid(rn *raft.RaftNode, vr raft.VoteResponse) bool {

	return false
}

func whatDoesTheMajorityWant(rn *raft.RaftNode, vr raft.VoteResponse) (bool, bool) {

	return false, false
}

func becomeAFollowerAccordingToLeader() {

}

func becomeALeader() {

}

func becomeAFollowerAccordingToPeer() {

}

func restartElection() {

}

func (em *ElectionManager) startElectionTimer(done chan bool, rn *raft.RaftNode) chan bool {
	electionTimedOut := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				break
			default:
				time.Sleep(em.ElectionTimeoutDuration)
				if rn.ElectionInProgress == true {
					rn.ElectionInProgress = false
				}

			}
		}
	}()
	return electionTimedOut
}

func (em *ElectionManager) askForVotes(done chan bool, rn *raft.RaftNode) chan raft.VoteResponse {
	votesReceived := make(chan raft.VoteResponse)
	for {
		select {
		case <-done:
			close(votesReceived)
			break
		default:
			voteRequest := em.GenerateVoteRequest(rn)
			for _, peer := range rn.Peers {
				go func(p raft.Peer) {
					votesReceived <- rn.RPCAdapter.RequestVoteFromPeer(p, voteRequest)
				}(peer)
			}
		}
	}

	return votesReceived
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
	eu := <-ec
	if !eu.ElectionTimerStarted {
		panic("Election Timer not started")
	}
}

func (em *ElectionManager) GetResponseForVoteRequest(rn *raft.RaftNode, vr raft.VoteRequest) raft.VoteResponse {
	return getVoteResponseForVoteRequest(rn, vr)
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
