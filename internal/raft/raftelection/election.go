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
	leaderHeartbeatChannel chan raft.Peer
)

type ElectionManager struct {
	ElectionTimeoutDuration time.Duration
	ElectionTimerOff        bool
}

func (em *ElectionManager) GetLeaderHeartChannel() chan raft.Peer {
	if leaderHeartbeatChannel == nil {
		leaderHeartbeatChannel = make(chan raft.Peer)
	}

	return leaderHeartbeatChannel
}

func (el *ElectionManager) HasElectionTimerStarted() bool {
	return !el.ElectionTimerOff
}

func (el *ElectionManager) HasElectionTimerStopped() bool {
	return el.ElectionTimerOff
}

func (em *ElectionManager) StartElection(rn *raft.RaftNode) {
	electionTimerDone := make(chan bool)
	votesReceived := make(chan raft.VoteResponse)
	conclusionDone := make(chan bool)

	rn.Mu.Lock()
	defer rn.Mu.Unlock()

	rn.LeaderHeartbeatMonitor.Stop()
	rn.ElectionInProgress = true
	rn.CurrentRole = raft.CANDIDATE
	rn.VotedFor = rn.Id
	rn.CurrentTerm = rn.CurrentTerm + 1
	rn.VotesReceived = nil

	em.askForVotes(votesReceived, rn)
	electionTimedOut := em.startElectionTimer(electionTimerDone, rn)
	followerAccordingToPeer, leader, follower := concludeFromReceivedVotes(conclusionDone, rn, votesReceived)

	for {
		select {
		case vr := <-followerAccordingToPeer:
			becomeAFollowerAccordingToPeer(rn, vr)
		case <-leader:
			becomeALeader(rn)
			break
		case l := <-leaderHeartbeatChannel:
			closeElectionChannels(electionTimerDone, votesReceived, conclusionDone)
			becomeAFollowerAccordingToLeader(l)
			break
		case <-follower:
			becomeAFollower(rn)
			closeElectionChannels(electionTimerDone, votesReceived, conclusionDone)
			break
		case <-electionTimedOut:
			closeElectionChannels(electionTimerDone, votesReceived, conclusionDone)
			restartElection()
			break
		}
	}
}

func closeElectionChannels(electionTimerDone chan bool, votesReceived chan raft.VoteResponse, conclusionDone chan bool) {
	close(electionTimerDone)
	close(votesReceived)
	close(conclusionDone)
}

func becomeAFollower(rn *raft.RaftNode) {
	rn.CurrentRole = raft.FOLLOWER
	rn.ElectionInProgress = false
	rn.IsHeartBeating = false
}

func becomeALeader(rn *raft.RaftNode) {
	rn.CurrentRole = raft.LEADER
	rn.ElectionInProgress = false
	rn.LogMgr.ReplicateLogs(rn)
	rn.Heart.StartBeating()
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
				votes := []raft.VoteResponse{}
				for vr := range votesReceived {
					votes = append(votes, vr)
					valid := isElectionValid(rn, vr)
					if !valid && isPeerTermHigher(rn, vr) {
						followerAccordingToPeer <- vr
					} else if becomeLeader, becomeFollower := whatDoesTheMajorityWant(len(rn.Peers), votes, vr); true {
						if becomeLeader {
							leader <- true
						} else if becomeFollower {
							follower <- true
						}
					}
				}
			}
		}
	}()

	return followerAccordingToPeer, leader, follower
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

func whatDoesTheMajorityWant(numOfPeers int, votesReceived []raft.VoteResponse, vr raft.VoteResponse) (bool, bool) {

	majorityCount := numOfPeers / 2
	votesInFavor := 0
	leader, follower := false, false
	for _, vr := range votesReceived {
		if vr.VoteGranted {
			votesInFavor = votesInFavor + 1
		}
	}

	// if majority has voted against, game over
	if len(votesReceived)-votesInFavor > majorityCount {
		leader = true
	}

	// if majority has voted in favor, game won
	if votesInFavor >= majorityCount {
		follower = false
	}

	return leader, follower
}

func becomeAFollowerAccordingToLeader(leader raft.Peer) {

}

func becomeAFollowerAccordingToPeer(rn *raft.RaftNode, v raft.VoteResponse) {
	rn.CurrentTerm = v.Term
	rn.CurrentRole = raft.FOLLOWER
	rn.VotedFor = 0
	rn.ElectionInProgress = false
	rn.IsHeartBeating = false
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
				electionTimedOut <- true
			}
		}
	}()
	return electionTimedOut
}

func (em *ElectionManager) askForVotes(votesReceived chan raft.VoteResponse, rn *raft.RaftNode) {
	voteRequest := em.GenerateVoteRequest(rn)
	for _, peer := range rn.Peers {
		go func(p raft.Peer) {
			votesReceived <- rn.RPCAdapter.RequestVoteFromPeer(p, voteRequest)
		}(peer)
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
